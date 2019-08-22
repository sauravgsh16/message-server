package server

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/sauravgsh16/secoc-third/qserver/consumer"
	"github.com/sauravgsh16/secoc-third/qserver/proto"
	"github.com/sauravgsh16/secoc-third/qserver/queue"
)

const (
	CH_INIT = iota
	CH_OPEN
	CH_CLOSING
	CH_CLOSED
)

type Channel struct {
	id             uint16
	server         *Server
	incoming       chan *proto.WireFrame
	outgoing       chan *proto.WireFrame
	conn           *Connection
	consumers      map[string]*consumer.Consumer
	consumerMux    sync.Mutex
	sendLock       sync.Mutex
	state          uint8
	currentMessage *proto.Message
	flow           bool
	usedQueueName  string
	deliveryTag    uint64
	tagMux         sync.Mutex
	txMode         bool
	txMessages     []*proto.TxMessage
	txLock         sync.Mutex
	defaultSize    uint32
	activeSize     uint32
	sizeMux        sync.Mutex
}

func NewChannel(id uint16, conn *Connection) *Channel {
	return &Channel{
		id:          id,
		server:      conn.server,
		incoming:    make(chan *proto.WireFrame),
		outgoing:    conn.outgoing,
		conn:        conn,
		consumers:   make(map[string]*consumer.Consumer),
		flow:        true,
		txMessages:  make([]*proto.TxMessage, 0),
		defaultSize: uint32(2048),
	}
}

func (ch *Channel) start() {
	if ch.state == 0 {
		ch.state = CH_OPEN
		go ch.startConnection()
	}

	go func() {
		for {
			if ch.state == CH_CLOSED {
				break
			}
			var err *proto.ProtoError
			frame := <-ch.incoming
			switch {
			case frame.FrameType == uint8(proto.FrameMethod):
				fmt.Println("routing method") // LOGS
				err = ch.routeMethod(frame)
			case frame.FrameType == uint8(proto.FrameHeader):
				if ch.state != CH_CLOSING {
					fmt.Println("handling header") // LOGS
					err = ch.handleHeader(frame)
				}
			case frame.FrameType == uint8(proto.FrameBody):
				if ch.state != CH_CLOSING {
					fmt.Println("handling body") // LOGS
					err = ch.handleBody(frame)
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

func (ch *Channel) SendMethod(m proto.MethodFrame) {
	buf := bytes.NewBuffer([]byte{})
	m.Write(buf)
	ch.outgoing <- &proto.WireFrame{
		FrameType: uint8(proto.FrameMethod),
		Channel:   ch.id,
		Payload:   buf.Bytes(),
	}
}

// **************IMPLEMENT BELOW ************************
// ********************************************************

func (ch *Channel) SendContent(mf proto.MethodFrame, msg *proto.Message) {}

//  *****************************************************
// *******************************************************

func (ch *Channel) FlowActive() bool {
	return ch.flow
}

func (ch *Channel) GetDeliveryTag() uint64 {
	ch.tagMux.Lock()
	defer ch.tagMux.Unlock()

	ch.deliveryTag++
	return ch.deliveryTag
}

func (ch *Channel) AcquireResources(qm *proto.QueueMessage) bool {
	ch.sizeMux.Lock()
	defer ch.sizeMux.Unlock()

	if ch.activeSize < ch.defaultSize {
		ch.activeSize += qm.MsgSize
		return true
	}
	return false
}

func (ch *Channel) ReleaseResources(qm *proto.QueueMessage) {
	ch.sizeMux.Lock()
	defer ch.sizeMux.Unlock()

	ch.activeSize -= qm.MsgSize
}

// ************************* PRIVATE METHODS ***********************************

func (ch *Channel) startTxMode() {
	ch.txMode = true
}

func (ch *Channel) commitTx(clsID, mtdID uint16) *proto.ProtoError {

	ch.txLock.Lock()
	defer ch.txLock.Unlock()

	mapQueueWithQueueMessages, err := ch.server.msgStore.AddTxMessages(ch.txMessages)
	if err != nil {
		return proto.NewSoftError(500, err.Error(), clsID, mtdID)
	}

	for qName, qMsgs := range mapQueueWithQueueMessages {
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

func (ch *Channel) rollbackTx() *proto.ProtoError {
	ch.txLock.Lock()
	defer ch.txLock.Unlock()

	ch.txMessages = make([]*proto.TxMessage, 0)
	return nil
}

func (ch *Channel) sendError(err *proto.ProtoError) {
	if err.Soft {
		fmt.Println("Sending channel error: ", err.Msg)
		ch.state = CH_CLOSING
		ch.SendMethod(&proto.ChannelClose{
			ReplyCode: err.Code,
			ReplyText: err.Msg,
			ClassId:   err.Class,
			MethodId:  err.Method,
		})
	} else {
		ch.conn.closeConnWithError(err)
	}
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

func (ch *Channel) addNewConsumer(q *queue.Queue, m *proto.BasicConsume) *proto.ProtoError {
	clsID, mtdID := m.MethodIdentifier()

	c := consumer.NewConsumer(ch.server.msgStore, ch, m.ConsumerTag, q, q.Name, m.NoAck)
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
	ch.currentMessage = proto.NewMessage(m)
}

func (ch *Channel) shutdown() {
	if ch.state == CH_CLOSED {
		fmt.Printf("channel already closed, shutdown performed on %d\n", ch.id)
		return
	}
	ch.state = CH_CLOSED
	// unregister channel from connection
	ch.conn.removeChannel(ch.id)
	// remove any consumer associated with this channel
	for _, c := range ch.consumers {
		ch.removeConsumer(c.ConsumerTag)
	}
}

func (ch *Channel) routeMethod(frame *proto.WireFrame) *proto.ProtoError {
	var methodReader = bytes.NewReader(frame.Payload)

	var methodFrame, err = proto.ReadMethod(methodReader)
	if err != nil {
		return proto.NewHardError(500, err.Error(), 0, 0)
	}

	var classID, methodID = methodFrame.MethodIdentifier()

	// Check if channel is in initial creation state
	if ch.state == CH_INIT && (classID != 20 || methodID != 10) {
		return proto.NewHardError(503, "Open method call on non-open channel", classID, methodID)
	}

	// Route methodFrame based on classID
	switch {
	case classID == 10:
		return ch.connectionRoute(ch.conn, methodFrame)
	case classID == 20:
		return ch.channelRoute(methodFrame)
	case classID == 30:
		return ch.exchangeRoute(methodFrame)
	case classID == 40:
		return ch.queueRoute(methodFrame)
	case classID == 50:
		return ch.basicRoute(methodFrame)
	case classID == 60:
		return ch.txRoute(methodFrame)
	default:
		return proto.NewHardError(540, "Not Implemented", classID, methodID)
	}
}

func (ch *Channel) handleHeader(frame *proto.WireFrame) *proto.ProtoError {

	if ch.currentMessage == nil {
		return proto.NewSoftError(500, "unexpected header frame", 0, 0)
	}

	if ch.currentMessage.Header != nil {
		return proto.NewSoftError(500, "unexpected - already seen header", 0, 0)
	}

	var header = &proto.HeaderFrame{}
	var buf = bytes.NewReader(frame.Payload)
	var err = header.Read(buf)
	if err != nil {
		return proto.NewHardError(500, "Error parsing header frame: "+err.Error(), 0, 0)
	}
	ch.currentMessage.Header = header
	return nil
}

func (ch *Channel) handleBody(frame *proto.WireFrame) *proto.ProtoError {

	if ch.currentMessage == nil {
		return proto.NewSoftError(500, "unexpected header frame", 0, 0)
	}

	if ch.currentMessage.Header == nil {
		return proto.NewSoftError(500, "unexpected body frame - no header yet", 0, 0)
	}

	ch.currentMessage.Payload = append(ch.currentMessage.Payload, frame)

	var size = uint64(0)
	for _, body := range ch.currentMessage.Payload {
		size += uint64(len(body.Payload))
	}
	// Message yet to complete
	if size < ch.currentMessage.Header.BodySize {
		return nil
	}

	ex, _ := ch.server.exchanges[ch.currentMessage.Method.Exchange]

	if ch.txMode {
		// Add message to a List
		queues, err := ex.QueuesToPublish(ch.currentMessage)
		if err != nil {
			return err
		}

		// Add TxMessage for all queues
		ch.txLock.Lock()
		for _, queueName := range queues {
			txMsg := proto.NewTxMessage(ch.currentMessage, queueName)
			ch.txMessages = append(ch.txMessages, txMsg)
		}
		ch.txLock.Unlock()
	} else {
		// Normal mode, publish directly
		returnMethod, err := ch.server.publish(ex, ch.currentMessage)
		if err != nil {
			ch.currentMessage = nil
			return err
		}
		if returnMethod != nil {
			ch.SendContent(returnMethod, ch.currentMessage)
		}
	}

	ch.currentMessage = nil
	return nil
}
