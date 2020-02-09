package server

import (
	"fmt"
	"sync"

	"github.com/sauravgsh16/message-server/proto"
	"github.com/sauravgsh16/message-server/qserver/consumer"
	"github.com/sauravgsh16/message-server/qserver/queue"
)

const (
	chInit = iota
	chOpen
	chClosing
	chClosed
)

// Channel struct
type Channel struct {
	id            uint16
	server        *Server
	incoming      chan proto.Frame
	outgoing      chan proto.Frame
	conn          *Connection
	consumers     map[string]*consumer.Consumer
	consumerMux   sync.Mutex
	sendMux       sync.Mutex
	state         uint8
	curMsg        *proto.Message
	flow          bool
	usedQueueName string
	deliveryTag   uint64
	tagMux        sync.Mutex
	txMode        bool
	txMessages    []*proto.TxMessage
	txLock        sync.Mutex
	defaultSize   uint32
	activeSize    uint32
	sizeMux       sync.Mutex
}

// NewChannel returns a new channel
func NewChannel(id uint16, conn *Connection) *Channel {
	return &Channel{
		id:          id,
		server:      conn.server,
		incoming:    make(chan proto.Frame),
		outgoing:    conn.outgoing,
		conn:        conn,
		consumers:   make(map[string]*consumer.Consumer),
		flow:        true,
		txMessages:  make([]*proto.TxMessage, 0),
		defaultSize: uint32(2048),
	}
}

// Send takes in a message frame and writes it on the connection
func (ch *Channel) Send(msgf proto.MessageFrame) error {

	fmt.Printf("Sending: %s\n", msgf.MethodName())

	if ch.state == chClosed {
		return ch.sendClosed(msgf)
	}

	return ch.sendOpen(msgf)
}

func (ch *Channel) sendClosed(msgf proto.MessageFrame) error {
	// We just need to send ChannelCloseOK
	if _, ok := msgf.(*proto.ChannelCloseOk); ok {
		return ch.conn.send(&proto.MethodFrame{
			ChannelID: ch.id,
			Method:    msgf,
		})
	}
	clsID, mtdID := msgf.Identifier()
	return proto.NewSoftError(501, "Send attempt on closed channel", clsID, mtdID)
}

func (ch *Channel) sendOpen(msgf proto.MessageFrame) error {

	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	if mcf, ok := msgf.(proto.MessageContentFrame); ok {

		prop, body := mcf.GetContent()
		clsID, _ := mcf.Identifier()
		size := uint64(len(body))

		// Send Method
		ch.outgoing <- &proto.MethodFrame{
			ChannelID: ch.id,
			Method:    mcf,
		}

		// Send Header
		ch.outgoing <- &proto.HeaderFrame{
			ChannelID:  ch.id,
			Class:      clsID,
			BodySize:   size,
			Properties: prop,
		}

		// Send Body
		ch.outgoing <- &proto.BodyFrame{
			ChannelID: ch.id,
			Body:      body,
		}
	} else {
		ch.outgoing <- &proto.MethodFrame{
			ChannelID: ch.id,
			Method:    msgf,
		}
	}
	return nil
}

// SendContent takes in a message content frame to write on the tcp connection
func (ch *Channel) SendContent(mcf proto.MessageContentFrame, msg *proto.Message) error {
	mcf.SetContent(msg.Header.Properties, msg.Payload)
	if err := ch.Send(mcf); err != nil {
		return err
	}
	return nil
}

// FlowActive returns true for active flow, false otherwise
func (ch *Channel) FlowActive() bool {
	return ch.flow
}

// GetDeliveryTag increments and returns the delivery tag for the message
func (ch *Channel) GetDeliveryTag() uint64 {
	ch.tagMux.Lock()
	defer ch.tagMux.Unlock()

	ch.deliveryTag++
	return ch.deliveryTag
}

// AcquireResources increments the active size of the message being sent if less than default allowed size
// Returns true if incremented, false otherwise
func (ch *Channel) AcquireResources(qm *proto.QueueMessage) bool {
	ch.sizeMux.Lock()
	defer ch.sizeMux.Unlock()

	if ch.activeSize < ch.defaultSize {
		ch.activeSize += qm.MsgSize
		return true
	}
	return false
}

// ReleaseResources decrements the active message size by the current message being sent
func (ch *Channel) ReleaseResources(qm *proto.QueueMessage) {
	ch.sizeMux.Lock()
	defer ch.sizeMux.Unlock()

	ch.activeSize -= qm.MsgSize
}

func (ch *Channel) start() {
	if ch.state == 0 && ch.id == 0 {
		ch.state = chOpen
		go ch.startConnection()
	}

	go func() {
		for {
			if ch.state == chClosed {
				break
			}
			var err *proto.Error
			frame := <-ch.incoming

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
				err = proto.NewHardError(500, "Unknown frame type: ", 0, 0)
			}
			if err != nil {
				ch.sendError(err)
			}
		}
	}()
}

func (ch *Channel) sendError(err *proto.Error) {
	if err.Soft {
		fmt.Println("Sending channel error: ", err.Msg)
		ch.state = chClosing
		ch.Send(&proto.ChannelClose{
			ReplyCode: err.Code,
			ReplyText: err.Msg,
			ClassId:   err.Class,
			MethodId:  err.Method,
		})
	} else {
		ch.conn.closeConnWithError(err)
	}
}

func (ch *Channel) shutdown() {
	if ch.state == chClosed {
		fmt.Printf("channel already closed, shutdown performed on %d\n", ch.id)
		return
	}
	ch.state = chClosed
	// unregister channel from connection
	ch.conn.removeChannel(ch.id)
	// remove any consumer associated with this channel
	for _, c := range ch.consumers {
		ch.removeConsumer(c.ConsumerTag)
	}
}

func (ch *Channel) close(code uint16, text string, clsID uint16, mtdID uint16) {
	ch.Send(&proto.ChannelClose{
		ReplyCode: code,
		ReplyText: text,
		ClassId:   clsID,
		MethodId:  mtdID,
	})
	ch.state = chClosing
}

func (ch *Channel) startTxMode() {
	ch.txMode = true
}

func (ch *Channel) commitTx(clsID, mtdID uint16) *proto.Error {

	ch.txLock.Lock()
	defer ch.txLock.Unlock()

	qQueueMsgMap, err := ch.server.msgStore.AddTxMessages(ch.txMessages)
	if err != nil {
		return proto.NewSoftError(500, err.Error(), clsID, mtdID)
	}

	for qName, qMsgs := range qQueueMsgMap {
		queue, found := ch.server.queues[qName]
		if !found {
			continue
		}
		for _, qMsg := range qMsgs {
			if !queue.Add(qMsg) {
				// If adding message to queue was not successful,
				// It means queue was closed.
				// We thus need to remove reference of the message store.
				resourceHolder := []proto.MessageResourceHolder{ch}
				ch.server.msgStore.RemoveRef(qMsg, qName, resourceHolder)
			}
		}
	}

	// Clear transaction
	ch.txMessages = make([]*proto.TxMessage, 0)

	return nil
}

func (ch *Channel) rollbackTx() *proto.Error {
	ch.txLock.Lock()
	defer ch.txLock.Unlock()

	ch.txMessages = make([]*proto.TxMessage, 0)
	return nil
}

func (ch *Channel) activateFlow(active bool) {
	if ch.flow == active {
		return
	}
	// Change flow to active
	ch.flow = active
	// Ping Consumers to start work again, if possible
	if ch.flow {
		for _, c := range ch.consumers {
			c.Ping()
		}
	}
}

func (ch *Channel) addNewConsumer(q *queue.Queue, m *proto.BasicConsume) *proto.Error {
	clsID, mtdID := m.Identifier()

	c := consumer.NewConsumer(ch.server.msgStore, ch, m.ConsumerTag, q, q.Name, m.NoAck, ch.defaultSize)
	ch.consumerMux.Lock()
	defer ch.consumerMux.Unlock()

	_, found := ch.consumers[c.ConsumerTag]
	if found {
		return proto.NewHardError(520, "Consumer already present", clsID, mtdID) // CHECK THE ERROR CODE -- STORE ALL ERROR CODE
	}

	// Add consumer to queue
	code, err := q.AddConsumer(c)
	if err != nil {
		return proto.NewSoftError(code, err.Error(), clsID, mtdID)
	}

	// Add consumer to the channel
	ch.consumers[c.ConsumerTag] = c

	c.Start()
	return nil
}

func (ch *Channel) removeConsumer(consumerTag string) error {
	c, found := ch.consumers[consumerTag]
	if !found {
		return fmt.Errorf("Consumer: %s not found", consumerTag)
	}

	c.Stop()
	ch.consumerMux.Lock()
	delete(ch.consumers, consumerTag)
	ch.consumerMux.Unlock()

	return nil
}

func (ch *Channel) startPublish(m *proto.BasicPublish) {
	ch.curMsg = proto.NewMessage(m)
}

func (ch *Channel) startConnection() *proto.Error {
	ch.Send(&proto.ConnectionStart{
		Version:    1,
		Mechanisms: "PLAIN",
	})
	return nil
}

func (ch *Channel) handleMethod(mf *proto.MethodFrame) *proto.Error {

	// Check if channel is in initial creation state
	if ch.state == chInit && (mf.ClassID != 20 || mf.MethodID != 10) {
		return proto.NewHardError(503, "Open method call on non-open channel", mf.ClassID, mf.MethodID)
	}

	fmt.Println("Received: ", mf.Method.MethodName())

	// Route methodFrame based on clsID
	switch mf.ClassID {
	case 10:
		return ch.connRoute(ch.conn, mf.Method)
	case 20:
		return ch.channelRoute(mf.Method)
	case 30:
		return ch.exchangeRoute(mf.Method)
	case 40:
		return ch.queueRoute(mf.Method)
	case 50:
		return ch.basicRoute(mf.Method)
	case 60:
		return ch.txRoute(mf.Method)
	default:
		return proto.NewHardError(540, "Not Implemented", mf.ClassID, mf.MethodID)
	}
}

func (ch *Channel) handleHeader(hf *proto.HeaderFrame) *proto.Error {
	if ch.curMsg == nil {
		return proto.NewSoftError(500, "unexpected header frame", 0, 0)
	}

	if ch.curMsg.Header != nil {
		return proto.NewSoftError(500, "unexpected - header already seen", 0, 0)
	}

	ch.curMsg.Header = hf

	return nil
}

func (ch *Channel) handleBody(bf *proto.BodyFrame) *proto.Error {
	if ch.curMsg == nil {
		return proto.NewSoftError(500, "unexpected header frame", 0, 0)
	}

	if ch.curMsg.Header == nil {
		return proto.NewSoftError(500, "unexpected body frame - no header yet", 0, 0)
	}

	ch.curMsg.Payload = append(ch.curMsg.Payload, bf.Body...)

	size := uint64(len(ch.curMsg.Payload))

	// Message yet to complete, we return
	if size < ch.curMsg.Header.BodySize {
		return nil
	}

	ex, _ := ch.server.getExchange(ch.curMsg.Method.(*proto.BasicPublish).Exchange)

	if ch.txMode {
		// Add message to a List
		queues, err := ex.QueuesToPublish(ch.curMsg)
		if err != nil {
			return err
		}

		// Add TxMessage for all queues
		ch.txLock.Lock()
		for _, queueName := range queues {
			txMsg := proto.NewTxMessage(ch.curMsg, queueName)
			ch.txMessages = append(ch.txMessages, txMsg)
		}
		ch.txLock.Unlock()
	} else {
		// Normal mode, publish directly
		returnMtd, err := ch.server.publish(ex, ch.curMsg)
		if err != nil {
			ch.curMsg = nil
			return err
		}
		if returnMtd != nil {
			ch.SendContent(returnMtd, ch.curMsg)
		}
	}

	ch.curMsg = nil

	return nil
}
