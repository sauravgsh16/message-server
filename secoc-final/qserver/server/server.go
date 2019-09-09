package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/boltdb/bolt"

	"github.com/sauravgsh16/secoc-third/proto"
	"github.com/sauravgsh16/secoc-third/qserver/exchange"
	"github.com/sauravgsh16/secoc-third/qserver/queue"
	"github.com/sauravgsh16/secoc-third/qserver/store"
)

type Server struct {
	exchanges       map[string]*exchange.Exchange
	queues          map[string]*queue.Queue
	conns           map[int64]*Connection
	mux             sync.Mutex
	db              *bolt.DB
	msgStore        *store.MsgStore
	exchangeDeleter chan *exchange.Exchange
	queueDeleter    chan *queue.Queue
}

// INCASE - THE SERVER AND THE MESSAGE DB NEEDS TO BE SEPARATE - THIS IS THE POINT
// WHERE WE ACCEPT TO DIFFERENT DB PATHS.
func NewServer(dbFilePath, msgStoreFilePath string) *Server {
	db, err := bolt.Open(dbFilePath, 0666, nil)
	if err != nil {
		panic(err.Error())
	}
	msgStore, err := store.New(msgStoreFilePath)
	if err != nil {
		panic("unable to create message store")
	}
	msgStore.Start()
	var s = &Server{
		exchanges:       make(map[string]*exchange.Exchange),
		queues:          make(map[string]*queue.Queue),
		conns:           make(map[int64]*Connection),
		exchangeDeleter: make(chan *exchange.Exchange),
		queueDeleter:    make(chan *queue.Queue),
		db:              db,
		msgStore:        msgStore,
	}

	s.initSystemExchanges()

	go s.deleteExchangeMonitor()
	go s.deleteQueueMonitor()
	return s
}

func (s *Server) OpenConnection(conn net.Conn) {
	c := NewConnection(s, conn)
	s.conns[c.id] = c
	c.openConnection()
}

// ****** PRIVATE METHODS *********

func (s *Server) initSystemExchanges() {
	s.registerDefaultExchange("proto.DIRECT", exchange.EX_DIRECT)
	s.registerDefaultExchange("proto.FANOUT", exchange.EX_FANOUT)
}

func (s *Server) registerDefaultExchange(name string, extype uint8) {
	_, alreadyPresent := s.exchanges[name]
	if !alreadyPresent {
		ex := exchange.NewExchange(
			name,
			extype,
			s.exchangeDeleter,
		)
		// TODO
		// PERSIST DB -- WHEN DB IS IMPLEMENTED
		s.addExchange(ex)
	}
}

func (s *Server) deleteExchangeMonitor() {
	for e := range s.exchangeDeleter {
		exDel := &proto.ExchangeDelete{
			Exchange: e.Name,
			NoWait:   true,
		}
		s.deleteExchange(exDel)
	}
}

func (s *Server) deleteQueueMonitor() {
	for q := range s.queueDeleter {
		qDel := &proto.QueueDelete{
			Queue:  q.Name,
			NoWait: true,
		}
		s.deleteQueue(qDel, -1)
	}
}

func (s *Server) addExchange(ex *exchange.Exchange) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.exchanges[ex.Name] = ex
	return nil
}

func (s *Server) deleteExchange(m *proto.ExchangeDelete) (uint16, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	ex, ok := s.exchanges[m.Exchange]
	if !ok {
		return 404, fmt.Errorf("Exchange: %s - not found", m.Exchange)
	}
	// Close everything associated with the exchange
	ex.Close()
	delete(s.exchanges, m.Exchange)
	return 0, nil
}

func (s *Server) addQueue(q *queue.Queue) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.queues[q.Name] = q
	return nil
}

func (s *Server) deleteQueue(m *proto.QueueDelete, connID int64) (uint32, uint16, error) {

	s.mux.Lock()
	defer s.mux.Unlock()

	q, found := s.queues[m.Queue]
	if !found {
		return 0, 404, fmt.Errorf("Queue not found")
	}

	if q.ConnId != -1 && q.ConnId != connID {
		return 0, 405, fmt.Errorf("Queue is locked by another connection")
	}

	// Close queue - to stop any data enqueue and dequeue
	q.Close()
	// Remove queue from all the bindings
	s.removeQueueBindings(m.Queue)

	// Cleanup
	msgPurged, err := q.Delete(m.IfUnused, m.IfEmpty)
	if err != nil {
		return 0, 406, err
	}
	delete(s.queues, m.Queue)
	return msgPurged, 0, nil
}

func (s *Server) removeQueueBindings(qName string) {
	for _, ex := range s.exchanges {
		ex.RemoveQueueBindings(qName)
	}
}

func (s *Server) deleteConnection(connID int64) {
	delete(s.conns, connID)
}

func (s *Server) deleteQueuesForConn(connID int64) {
	s.mux.Lock()
	qToDelete := make([]*queue.Queue, 0)
	for _, q := range s.queues {
		if q.ConnId == connID {
			qToDelete = append(qToDelete, q)
		}
	}
	s.mux.Unlock()

	for _, q := range qToDelete {
		qd := &proto.QueueDelete{
			Queue: q.Name,
		}
		s.deleteQueue(qd, connID)
	}
}

func (s *Server) basicReturnMsg(msg *proto.Message, code uint16, text string) *proto.BasicReturn {
	return &proto.BasicReturn{
		ReplyCode:  code,
		ReplyText:  text,
		Exchange:   msg.Method.(*proto.BasicReturn).Exchange,
		RoutingKey: msg.Method.(*proto.BasicReturn).RoutingKey,
	}
}

func (s *Server) publish(ex *exchange.Exchange, msg *proto.Message) (*proto.BasicReturn, *proto.Error) {
	if ex.Closed {
		return s.basicReturnMsg(msg, 313, "Exchange closed, unable to route message"), nil // AGAIN CHECK FOR RETURN CODE - IMPLEMENT CONSTANT
	}

	queues, err := ex.QueuesToPublish(msg)
	if err != nil {
		return nil, err
	}

	// No avaliable queues
	if len(queues) == 0 {
		return s.basicReturnMsg(msg, 313, "No available queues found"), nil
	}

	// Add message and queue to message store.
	mapQueueWithQueueMessages, errObj := s.msgStore.AddMessage(msg, queues)
	if errObj != nil {
		clsID, mtdID := msg.Method.MethodIdentifier()
		return nil, proto.NewSoftError(500, errObj.Error(), clsID, mtdID)
	}

	if msg.Method.(*proto.BasicPublish).Immediate {
		return s.consumeMsgImmediate(msg, queues, mapQueueWithQueueMessages)
	}
	s.addMsgForConsumption(msg, queues, mapQueueWithQueueMessages)
	return nil, nil
}

func (s *Server) consumeMsgImmediate(msg *proto.Message, queues []string, qmMap map[string][]*proto.QueueMessage) (*proto.BasicReturn, *proto.Error) {
	consumed := false
	for _, queueName := range queues {
		qms := qmMap[queueName]
		for _, qm := range qms {
			queue, found := s.queues[queueName]
			if !found {
				// Queue could have been deleted
				continue
			}
			msgConsumed := queue.ConsumeImmediate(qm)
			mrh := make([]proto.MessageResourceHolder, 0)
			if !msgConsumed {
				s.msgStore.RemoveRef(qm, queueName, mrh)
			}
			consumed = consumed || msgConsumed
		}
	}
	if !consumed {
		return s.basicReturnMsg(msg, 313, "No consumers available"), nil
	}
	return nil, nil
}

func (s *Server) addMsgForConsumption(msg *proto.Message, queues []string, qmMap map[string][]*proto.QueueMessage) {
	for _, queueName := range queues {
		qMsgs := qmMap[queueName]
		for _, qm := range qMsgs {
			q, found := s.queues[queueName]
			if !found || !q.Add(qm) {
				// Need to remove queue reference from msg store
				// particular queue message
				mrh := make([]proto.MessageResourceHolder, 0)
				s.msgStore.RemoveRef(qm, queueName, mrh)
			}
		}
	}
}