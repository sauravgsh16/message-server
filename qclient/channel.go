package qclient

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/sauravgsh16/message-server/proto"
)

const (
	// Channel status
	chInit = iota
	chOpen
	chClosing
	chClosed
)

// TODO: MOVE TO RELEVANT PLACE

// Return struct
type Return struct{}

// Confirmation struct
type Confirmation struct {
	DeliveryTag uint64
	State       bool
}

type confirms struct{}

func (c *confirms) AddListener(ch chan Confirmation) {}

// TODO: to move code ends here

// MetaData struct used when publishing message. Describe the metadata of the message
type MetaDataWithBody struct {
	ContentType   string
	MessageID     string
	UserID        string
	ApplicationID string
	Body          []byte
}

// Channel struct
type Channel struct {
	id              uint16
	destructor      sync.Once
	incoming        chan proto.Frame
	outgoing        chan proto.Frame
	outgoingContent chan proto.Frame
	rpc             chan proto.MessageFrame
	conn            *Connection
	consumers       *Consumers
	sendMux         sync.Mutex
	notifyMux       sync.Mutex
	state           uint8
	errors          chan *proto.Error
	confirms        *confirms
	flows           []chan bool
	cancels         []chan string
	closes          []chan *proto.Error
	returns         []chan Return
	noNotify        bool
	currentMsg      *proto.Message
	bodyMf          proto.MessageContentFrame
	done            chan interface{}
	contentWg       *sync.WaitGroup
}

func newChannel(c *Connection, id uint16, wg *sync.WaitGroup) *Channel {
	return &Channel{
		id:              id,
		conn:            c,
		incoming:        make(chan proto.Frame),
		outgoing:        c.outgoing,
		outgoingContent: c.outgoingContent,
		rpc:             make(chan proto.MessageFrame),
		consumers:       CreateNewConsumers(),
		done:            make(chan interface{}),
		errors:          make(chan *proto.Error),
		contentWg:       wg,
	}
}

func (ch *Channel) call(req proto.MessageFrame, resp ...proto.MessageFrame) error {
	if err := ch.send(req); err != nil {
		return err
	}

	if req.Wait() {
		select {
		case e, ok := <-ch.errors:
			if ok {
				return e
			}
			return ErrClosed

		case msg := <-ch.rpc:
			if msg != nil {
				for _, res := range resp {
					if reflect.TypeOf(res) == reflect.TypeOf(msg) {
						vres := reflect.ValueOf(res).Elem()
						vmsg := reflect.ValueOf(msg).Elem()
						vres.Set(vmsg)
						return nil
					}
				}
				return ErrInvalidCommand
			}
			return ErrClosed
		}
	}
	return nil
}

func (ch *Channel) send(msgf proto.MessageFrame) error {

	fmt.Printf("Sending: %s\n", msgf.MethodName())

	if ch.state == chClosed {
		return ch.sendClosed(msgf)
	}

	return ch.sendOpen(msgf)
}

func (ch *Channel) sendOpen(msgf proto.MessageFrame) error {

	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	if mcf, ok := msgf.(proto.MessageContentFrame); ok {

		prop, body := mcf.GetContent()
		clsID, _ := mcf.Identifier()
		size := uint64(len(body))

		ch.contentWg.Add(1)

		// Send Method
		ch.outgoingContent <- &proto.MethodFrame{
			ChannelID: ch.id,
			Method:    mcf,
		}

		// Send Header
		ch.outgoingContent <- &proto.HeaderFrame{
			Class:      clsID,
			ChannelID:  ch.id,
			BodySize:   size,
			Properties: prop,
		}

		// Send Body
		ch.outgoingContent <- &proto.BodyFrame{
			ChannelID: ch.id,
			Body:      body,
		}

		ch.contentWg.Wait()

	} else {
		ch.outgoing <- &proto.MethodFrame{
			ChannelID: ch.id,
			Method:    msgf,
		}
	}
	return nil
}

func (ch *Channel) sendClosed(msgf proto.MessageFrame) error {
	if _, ok := msgf.(*proto.ChannelCloseOk); ok {
		ch.conn.send(&proto.MethodFrame{
			ChannelID: ch.id,
			Method:    msgf,
		})
	}
	return ErrClosed
}

func (ch *Channel) sendError(err *proto.Error) {
	if err.Soft {
		ch.send(&proto.ChannelClose{
			ReplyCode: err.Code,
			ReplyText: err.Msg,
			ClassId:   err.Class,
			MethodId:  err.Method,
		})
	} else {
		ch.conn.closeWithErr(err)
	}
}

func (ch *Channel) open() error {
	ch.startReceiver()
	return ch.call(&proto.ChannelOpen{}, &proto.ChannelOpenOk{Response: "200"})
}

func (ch *Channel) shutdown(err *proto.Error) {
	ch.destructor.Do(func() {
		ch.sendMux.Lock()
		defer ch.sendMux.Unlock()

		ch.notifyMux.Lock()
		defer ch.notifyMux.Unlock()

		if err != nil {
			for _, c := range ch.closes {
				c <- err
			}
		}

		ch.state = chClosed

		// Notify select loop for ch.rpc
		if err != nil {
			ch.errors <- err
		}

		ch.consumers.close()

		for _, c := range ch.closes {
			close(c)
		}

		for _, f := range ch.flows {
			close(f)
		}

		for _, r := range ch.returns {
			close(r)
		}

		for _, ca := range ch.cancels {
			close(ca)
		}

		ch.closes = nil
		ch.flows = nil
		ch.returns = nil
		ch.cancels = nil

		close(ch.errors)
	})
}

func (ch *Channel) startReceiver() {
	if ch.state == 0 {
		ch.state = chOpen
	}
	go func() {
		for {
			if ch.state == chClosed {
				break
			}
			var err *proto.Error

			select {
			case <-ch.done:
				break

			case frame := <-ch.incoming:

				switch m := frame.(type) {

				case *proto.MethodFrame:
					err = ch.handleMethod(m)

				case *proto.HeaderFrame:
					if ch.state != chClosing {
						err = ch.handleHeader(m)
					}

				case *proto.BodyFrame:
					if ch.state != chClosing {
						err = ch.handleBody(m)
					}

				default:
					err = proto.NewHardError(500, "Unknown frame type", 0, 0)
				}
				if err != nil {
					ch.sendError(err)
				}
			}
		}
	}()
}

func (ch *Channel) dispatchRPC(msgf proto.MessageFrame) *proto.Error {

	switch m := msgf.(type) {

	case *proto.ChannelClose:
		ch.sendMux.Lock()
		ch.send(&proto.ChannelCloseOk{})
		ch.sendMux.Unlock()
		ch.conn.closeChannel(ch, proto.NewSoftError(m.ReplyCode, m.ReplyText, m.ClassId, m.MethodId))

	case *proto.ChannelFlow:
		ch.notifyMux.Lock()
		for _, c := range ch.flows {
			c <- m.Active
		}
		ch.notifyMux.Unlock()
		ch.send(&proto.ChannelFlowOk{Active: m.Active})

	case *proto.BasicCancel:
		ch.notifyMux.Lock()
		for _, c := range ch.cancels {
			c <- m.ConsumerTag
		}
		ch.notifyMux.Unlock()
		ch.consumers.cancel(m.ConsumerTag)

	case *proto.BasicReturn:
		p := fmt.Sprintf("Not implemented %#v\n", m)
		panic(p)

	case *proto.BasicAck:
		panic("Not implemented")

	case *proto.BasicNack:
		panic("Not implemented")

	case *proto.BasicDeliver:
		ch.consumers.send(m.ConsumerTag, newDelivery(ch, m))

	default:
		ch.rpc <- msgf
	}

	return nil
}

func (ch *Channel) handleMethod(mf *proto.MethodFrame) *proto.Error {

	if msgf, ok := mf.Method.(proto.MessageContentFrame); ok {
		ch.currentMsg = proto.NewMessage(msgf)
		ch.bodyMf = msgf
		return nil
	}

	ch.dispatchRPC(mf.Method)
	return nil
}

func (ch *Channel) handleHeader(hf *proto.HeaderFrame) *proto.Error {
	if ch.currentMsg == nil {
		return proto.NewSoftError(500, "Unexpected header frame. No current message set", 0, 0)
	}

	if ch.currentMsg.Header != nil {
		return proto.NewSoftError(500, "Already seen header", 0, 0)
	}

	ch.currentMsg.Header = hf
	return nil
}

func (ch *Channel) handleBody(bf *proto.BodyFrame) *proto.Error {
	if ch.currentMsg == nil {
		return proto.NewSoftError(500, "Unexpected Body frame. No current message set", 0, 0)
	}

	if ch.currentMsg.Header == nil {
		return proto.NewSoftError(500, "Unexpected body frame. Header not set", 0, 0)
	}

	ch.currentMsg.Payload = append(ch.currentMsg.Payload, bf.Body...)

	size := uint64(len(ch.currentMsg.Payload))
	// Message yet to complete
	if size < ch.currentMsg.Header.BodySize {
		return nil
	}

	// Set MessageFrame's body with all content recieved
	ch.bodyMf.SetContent(ch.currentMsg.Header.Properties, ch.currentMsg.Payload)

	var err *proto.Error

	if err := ch.dispatchRPC(ch.bodyMf); err != nil {
		err = proto.NewSoftError(500, "Unable to dispatch method content frame", 0, 0)
	}

	ch.resetCurMsg()
	return err
}

func (ch *Channel) resetCurMsg() {
	ch.bodyMf = nil
	ch.currentMsg = nil
}

// Close signals done channel to close. Also closes connection
func (ch *Channel) Close() error {
	err := ch.call(
		&proto.ChannelClose{ReplyCode: 200},
		&proto.ChannelCloseOk{},
	)
	ch.done <- struct{}{}
	ch.conn.closeChannel(ch, nil)
	return err
}

// Flow calls ChannelFlow and exprets ChannelFlowOk
func (ch *Channel) Flow(active bool) error {
	return ch.call(
		&proto.ChannelFlow{Active: active},
		&proto.ChannelFlowOk{},
	)
}

// NotifyClose appends notification channel to closing list
func (ch *Channel) NotifyClose(c chan *proto.Error) chan *proto.Error {
	ch.notifyMux.Lock()
	defer ch.notifyMux.Unlock()

	if ch.noNotify {
		close(c)
	} else {
		ch.closes = append(ch.closes, c)
	}
	return c
}

// NotifyReturn appends return channel to returns list
func (ch *Channel) NotifyReturn(c chan Return) chan Return {
	ch.notifyMux.Lock()
	defer ch.notifyMux.Unlock()

	if ch.noNotify {
		close(c)
	} else {
		ch.returns = append(ch.returns, c)
	}
	return c
}

// NotifyFlow appends flow channel to flow list
func (ch *Channel) NotifyFlow(c chan bool) chan bool {
	ch.notifyMux.Lock()
	defer ch.notifyMux.Unlock()

	if ch.noNotify {
		close(c)
	} else {
		ch.flows = append(ch.flows, c)
	}
	return c
}

// NotifyCancel appends cancel channel to cancels list
func (ch *Channel) NotifyCancel(c chan string) chan string {
	ch.notifyMux.Lock()
	defer ch.notifyMux.Unlock()

	if ch.noNotify {
		close(c)
	} else {
		ch.cancels = append(ch.cancels, c)
	}
	return c
}

// NotifyConfirm confirms the message delivery
func (ch *Channel) NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64) {
	confirm := ch.NotifyPublish(make(chan Confirmation, len(ack)+len(nack)))

	go func() {
		for c := range confirm {
			if c.State == true {
				ack <- c.DeliveryTag
			} else {
				nack <- c.DeliveryTag
			}
		}
		close(ack)
		if ack != nack {
			close(nack)
		}
	}()
	return ack, nack
}

// NotifyPublish confirms message publish
func (ch *Channel) NotifyPublish(c chan Confirmation) chan Confirmation {
	ch.notifyMux.Lock()
	defer ch.notifyMux.Unlock()

	if ch.noNotify {
		close(c)
	} else {
		ch.confirms.AddListener(c)
	}
	return c
}

// ExchangeDeclare declares an exchange
func (ch *Channel) ExchangeDeclare(name, etype string, noWait bool) error {
	return ch.call(
		&proto.ExchangeDeclare{
			Exchange: name,
			Type:     etype,
			NoWait:   noWait,
		},
		&proto.ExchangeDeclareOk{},
	)
}

// ExchangeBind binds an exchange to a routing key
func (ch *Channel) ExchangeBind(dest, src, routingKey string, noWait bool) error {
	return ch.call(
		&proto.ExchangeBind{
			Destination: dest,
			Source:      src,
			RoutingKey:  routingKey,
			NoWait:      noWait,
		},
		&proto.ExchangeBindOk{},
	)
}

// ExchangeUnbind unbinds an exchange
func (ch *Channel) ExchangeUnbind(dest, src, routingKey string, noWait bool) error {
	return ch.call(
		&proto.ExchangeUnbind{
			Destination: dest,
			Source:      src,
			RoutingKey:  routingKey,
			NoWait:      noWait,
		},
		&proto.ExchangeUnbindOk{},
	)
}

// ExchangeDelete deletes an exchange
func (ch *Channel) ExchangeDelete(name string, ifunused, noWait bool) error {
	return ch.call(
		&proto.ExchangeDelete{
			Exchange: name,
			IfUnused: ifunused,
			NoWait:   noWait,
		},
		&proto.ExchangeDeleteOk{},
	)
}

// QueueDeclare declares a queue
func (ch *Channel) QueueDeclare(name string, noWait bool) (*proto.QueueDeclareOk, error) {
	req := &proto.QueueDeclare{
		Queue:  name,
		NoWait: noWait,
	}
	resp := &proto.QueueDeclareOk{}

	if err := ch.call(req, resp); err != nil {
		return &proto.QueueDeclareOk{}, err
	}
	if req.Wait() {
		return resp, nil
	}
	return &proto.QueueDeclareOk{Queue: name}, nil
}

// QueueBind binds a queue
func (ch *Channel) QueueBind(name, exchange, key string, noWait bool) error {
	return ch.call(
		&proto.QueueBind{
			Queue:      name,
			Exchange:   exchange,
			RoutingKey: key,
			NoWait:     noWait,
		},
		&proto.QueueBindOk{},
	)
}

// QueueUnbind unbinds queue
func (ch *Channel) QueueUnbind(name, exchange, key string) error {
	return ch.call(
		&proto.QueueUnbind{
			Queue:      name,
			Exchange:   exchange,
			RoutingKey: key,
		},
		&proto.QueueUnbindOk{},
	)
}

// QueueDelete deletes queue
func (ch *Channel) QueueDelete(name string, ifunused, ifempty, noWait bool) (int, error) {
	req := &proto.QueueDelete{
		Queue:    name,
		IfUnused: ifunused,
		IfEmpty:  ifempty,
		NoWait:   noWait,
	}
	resp := &proto.QueueDeleteOk{}
	return int(resp.MessageCnt), ch.call(req, resp)
}

// Publish a message
func (ch *Channel) Publish(exchange, key string, immediate bool, meta MetaDataWithBody) error {
	bp := &proto.BasicPublish{
		Exchange:   exchange,
		RoutingKey: key,
		Immediate:  immediate,
		Body:       meta.Body,
		Properties: proto.Properties{
			ContentType:   meta.ContentType,
			MessageID:     meta.MessageID,
			UserID:        meta.UserID,
			ApplicationID: meta.ApplicationID,
		},
	}
	ch.currentMsg = proto.NewMessage(bp)
	if err := ch.send(bp); err != nil {
		return err
	}
	ch.currentMsg = nil
	return nil
}

// Cancel a consumer
func (ch *Channel) Cancel(tag string, noWait bool) error {
	req := &proto.BasicCancel{
		ConsumerTag: tag,
		NoWait:      noWait,
	}
	resp := &proto.BasicCancelOk{}
	if err := ch.call(req, resp); err != nil {
		return err
	}

	ch.consumers.cancel(tag)
	return nil
}

// Consume messages
func (ch *Channel) Consume(queue, consumer string, noAck, noWait bool) (<-chan Delivery, error) {
	req := &proto.BasicConsume{
		Queue:       queue,
		ConsumerTag: consumer,
		NoAck:       noAck,
		NoWait:      noWait,
	}
	resp := &proto.BasicConsumeOk{}

	dChan := make(chan Delivery)
	ch.consumers.add(consumer, dChan)

	if err := ch.call(req, resp); err != nil {
		ch.consumers.cancel(consumer)
		return nil, err
	}

	return dChan, nil
}

// Ack message
func (ch *Channel) Ack(tag uint64, multiple bool) error {
	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	return ch.send(&proto.BasicAck{
		DeliveryTag: tag,
		Multiple:    multiple,
	})
}

// Nack not ack
func (ch *Channel) Nack(tag uint64, multiple bool, requeue bool) error {
	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	return ch.send(&proto.BasicNack{
		DeliveryTag: tag,
		Multiple:    multiple,
		Requeue:     requeue,
	})
}

// TxSelect transaction select
func (ch *Channel) TxSelect() error {
	return ch.call(
		&proto.TxSelect{},
		&proto.TxSelectOk{},
	)
}

// TxCommit transaction commit
func (ch *Channel) TxCommit() error {
	return ch.call(
		&proto.TxCommit{},
		&proto.TxCommitOk{},
	)
}

// TxRollBack transaction rollback
func (ch *Channel) TxRollBack() error {
	return ch.call(
		&proto.TxRollback{},
		&proto.TxRollbackOk{},
	)
}
