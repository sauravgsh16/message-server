package qclient

import (
	"bytes"
	"errors"
	"reflect"
	"sync"

	"github.com/sauravgsh16/secoc-third/proto"
)

const (
	CH_INIT = iota
	CH_OPEN
	CH_CLOSING
	CH_CLOSED
)

// TO MOVE TO RELEVANT PLACE
// ****************************************
type Return struct{}

type Confirmation struct {
	DeliveryTag uint64
	State       bool
}

type confirms struct{}

func (c *confirms) AddListener(ch chan Confirmation) {}

// *****************************************

type Channel struct {
	id             uint16
	destructor     sync.Once
	incoming       chan *proto.WireFrame
	outgoing       chan *proto.WireFrame
	rpc            chan *proto.MethodFrame
	conn           *Connection
	consumers      *Consumers
	sendMux        sync.Mutex
	notifyMux      sync.Mutex
	state          uint8
	errors         chan *proto.Error
	confirms       *confirms
	flows          []chan bool
	cancels        []chan string
	closes         []chan *proto.Error
	returns        []chan Return
	noNotify       bool
	body           []byte
	currentMessage *proto.Message
	bodyMf         proto.MethodContentFrame 
}

func newChannel(c *Connection, id uint16) *Channel {
	return &Channel{
		id:        id,
		conn:      c,
		incoming:  make(chan *proto.WireFrame),
		outgoing:  c.outgoing,
		rpc:       make(chan *proto.MethodFrame),
		consumers: CreateNewConsumers(),
	}
}

func (ch *Channel) call(req proto.MethodFrame, resp ...proto.MethodFrame) error {
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

func (ch *Channel) send(mf proto.MethodFrame) error {
	if ch.state == CH_CLOSED {
		return ch.sendClosed(mf)
	}
	if content, ok := mf.(proto.MethodContentFrame); ok {
		if err := ch.updateCurrentMessage(content); err != nil {
			return errors.New(err.Error())
		}
		ch.sendContent(ch.currentMessage, mf)
	}
	return ch.sendMethod(mf)
}

func (ch *Channel) updateCurrentMessage(mcf proto.MethodContentFrame) *proto.Error {
	clsID, mtdID := mcf.MethodIdentifier()
	body := mcf.GetBody()
	if ch.currentMessage == nil {
		return proto.NewSoftError(500, "Unexpected Message update", clsID, mtdID)
	}
	ch.currentMessage.Header = &proto.HeaderFrame{
		Class:    clsID,
		BodySize: uint64(len(body)),
	}

	// TODO:
	// BODY TO BE CHUNKED INTO AGREED MAX FRAME SIZE
	// SINCE CONFIGURATION HANDSHAKE IS NOT YET IMPLEMENTED, FRAME SIZE WILL
	// BE THE SIZE OF BODY.
	ch.currentMessage.Payload = append(ch.currentMessage.Payload,
		&proto.WireFrame{
			FrameType: proto.FrameBody,
			Channel:   ch.id,
			Payload:   body,
		})
	return nil
}

func (ch *Channel) sendClosed(mf proto.MethodFrame) error {
	if _, ok := mf.(*proto.ChannelCloseOk); ok {
		ch.conn.send(&proto.ChannelFrame{
			ChannelID: ch.id,
			Method:    mf,
		})
	}
	return ErrClosed
}

func (ch *Channel) sendMethod(mf proto.MethodFrame) error {
	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	if err := ch.conn.send(&proto.ChannelFrame{
		ChannelID: ch.id,
		Method:    mf,
	}); err != nil {
		return err
	}
	return nil
}

func (ch *Channel) sendContent(msg *proto.Message, mf proto.MethodFrame) error {
	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	buf := bytes.NewBuffer(make([]byte, 0))

	// Write Headers
	proto.WriteShort(buf, msg.Header.Class)
	proto.WriteLongLong(buf, msg.Header.BodySize)

	// Send method
	if err := ch.sendMethod(mf); err != nil {
		return err
	}

	// Send Header
	ch.outgoing <- &proto.WireFrame{
		FrameType: uint8(proto.FrameHeader),
		Channel:   ch.id,
		Payload:   buf.Bytes(),
	}

	// Send Body
	for _, body := range msg.Payload {
		body.Channel = ch.id
		ch.outgoing <- body
	}
	return nil
}

func (ch *Channel) sendError(err *proto.Error) {
	if err.Soft {
		ch.sendMethod(&proto.ChannelClose{
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
	return ch.call(&proto.ChannelOpen{}, &proto.ChannelOpenOk{})
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

		ch.state = CH_CLOSED

		// Notify select loop for ch.rpc
		if err != nil {
			ch.errors <- err
		}

		ch.consumers.close()

		for _, c := range ch.closes {
			close(c)
		}

		for _, c := range ch.flows {
			close(c)
		}

		for _, c := range ch.returns {
			close(c)
		}

		for _, c := range ch.cancels {
			close(c)
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
		ch.state = CH_OPEN
	}
	go func() {
		for {
			if ch.state == CH_CLOSED {
				break
			}
			var err *proto.Error
			wf := <-ch.incoming

			switch wf.FrameType {
			case uint8(proto.FrameMethod):
				err = ch.handleMethod(wf)

			case uint8(proto.FrameHeader):
				if ch.state != CH_CLOSING {
					err = ch.handleHeader(wf)
				}

			case uint8(proto.FrameBody):
				if ch.state != CH_CLOSING {
					err = ch.handleBody(wf)
				}

			default:
				err = proto.NewHardError(500, "Unknown frame type", 0, 0)
			}
			if err != nil {
				ch.sendError(err)
			}
		}
	}()
}

func (ch *Channel) dispatchRpc(mf proto.MethodFrame) *proto.Error {

	switch m := mf.(type) {

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
		panic("Not implemented")

	case *proto.BasicAck:
		panic("Not implemented")

	case *proto.BasicNack:
		panic("Not implemented")

	case *proto.BasicDeliver:
		ch.consumers.send(m.ConsumerTag, newDelivery(ch, m))

	default:
		ch.rpc <- &mf
	}

	return nil
}

func (ch *Channel) handleMethod(wf *proto.WireFrame) *proto.Error {
	r := bytes.NewReader(wf.Payload)

	mf, err := proto.ReadMethod(r)
	if err != nil {
		pErr := proto.NewHardError(500, err.Error(), 0, 0)
		return pErr
	}

	if content, ok := mf.(proto.MethodContentFrame); ok {
		ch.body = make([]byte, 0)
		ch.currentMessage = proto.NewMessage(content)
		ch.bodyMf = content
	} else {
		ch.dispatchRpc(mf)
	}
	return nil
}

func (ch *Channel) handleHeader(wf *proto.WireFrame) *proto.Error {
	if ch.currentMessage == nil {
		return proto.NewSoftError(500, "Unexpected header frame. No current message set", 0, 0)
	}

	if ch.currentMessage.Header != nil {
		return proto.NewSoftError(500, "Already seen header", 0, 0)
	}

	hf := &proto.HeaderFrame{}
	if err := hf.Read(bytes.NewReader(wf.Payload)); err != nil {
		return proto.NewHardError(500, "Error reading header", 0, 0)
	}

	ch.currentMessage.Header = hf
	return nil
}

func (ch *Channel) handleBody(wf *proto.WireFrame) *proto.Error {
	if ch.currentMessage == nil {
		return proto.NewSoftError(500, "Unexpected Body frame. No current message set", 0, 0)
	}

	if ch.currentMessage.Header == nil {
		return proto.NewSoftError(500, "Unexpected body frame. Header not set", 0, 0)
	}

	ch.currentMessage.Payload = append(ch.currentMessage.Payload, wf)

	size := uint64(0)
	for _, b := range ch.currentMessage.Payload {
		size += uint64(len(b.Payload))
	}

	if size < ch.currentMessage.Header.BodySize {
		// We still need to wait for additional bodyframes
		return nil
	}

	for _, wf := range ch.currentMessage.Payload {
		ch.body = append(ch.body, wf.Payload...)
	}

	var err *proto.Error

	ch.bodyMf.SetBody(ch.body)
	if err := ch.dispatchRpc(ch.bodyMf); err != nil {
		err = proto.NewSoftError(500, "Unable to dispatch method content frame", 0, 0)
	}

	ch.resetMessages()
	return err
}

func (ch *Channel) resetMessages() {
	ch.body = nil
	ch.bodyMf = nil
	ch.currentMessage = nil
}

func (ch *Channel) Close() error {
	defer ch.conn.closeChannel(ch, nil)
	return ch.call(
		&proto.ChannelClose{ReplyCode: 200},
		&proto.ChannelCloseOk{},
	)
}

func (ch *Channel) Flow(active bool) error {
	return ch.call(
		&proto.ChannelFlow{Active: active},
		&proto.ChannelFlowOk{},
	)
}

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

func (ch *Channel) Publish(exchange, key string, immediate bool, body []byte) error {
	bp := &proto.BasicPublish{
		Exchange:   exchange,
		RoutingKey: key,
		Immediate:  immediate,
		Body:       body,
	}
	ch.currentMessage = proto.NewMessage(bp)
	if err := ch.send(bp); err != nil {
		return err
	}
	ch.currentMessage = nil
	return nil
}

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

func (ch *Channel) Ack(tag uint64, multiple bool) error {
	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	return ch.send(&proto.BasicAck{
		DeliveryTag: tag,
		Multiple:    multiple,
	})
}

func (ch *Channel) Nack(tag uint64, multiple bool, requeue bool) error {
	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	return ch.send(&proto.BasicNack{
		DeliveryTag: tag,
		Multiple:    multiple,
		Requeue:     requeue,
	})
}

func (ch *Channel) TxSelect() error {
	return ch.call(
		&proto.TxSelect{},
		&proto.TxSelectOk{},
	)
}

func (ch *Channel) TxCommit() error {
	return ch.call(
		&proto.TxCommit{},
		&proto.TxCommitOk{},
	)
}

func (ch *Channel) TxRollBack() error {
	return ch.call(
		&proto.TxRollback{},
		&proto.TxRollbackOk{},
	)
}
