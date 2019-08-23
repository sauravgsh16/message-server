package exchange

import (
	"fmt"
	"sync"

	"github.com/sauravgsh16/secoc-third/qserver/binding"
	"github.com/sauravgsh16/secoc-third/qserver/proto"
)

const (
	EX_DIRECT  uint8 = 1
	EX_FANOUT  uint8 = 2
	EX_HEADERS uint8 = 3
)

type Exchange struct {
	Name       string
	ExType     uint8
	bindings   []*binding.Binding
	bindLock   sync.Mutex
	incoming   chan proto.Frame
	Closed     bool
	deleteChan chan *Exchange
}

func NewExchange(name string, extype uint8, deleteChan chan *Exchange) *Exchange {
	return &Exchange{
		Name:       name,
		ExType:     extype,
		bindings:   make([]*binding.Binding, 0),
		incoming:   make(chan proto.Frame),
		deleteChan: deleteChan,
	}
}

// NEED TO CHECK - RABBITMQ DOC
// GETS CALLED FROM exchangeDeclare
func NewExchangeFromMethod(m *proto.ExchangeDeclare, exDeleter chan *Exchange) (*Exchange, *proto.ProtoError) {
	extype, err := exchangeNameToType(m.Type)
	if err != nil {
		var classId, methodId = m.MethodIdentifier()
		return nil, proto.NewHardError(503, "Invalid exchange type", classId, methodId)
	}

	ex := NewExchange(m.Exchange, extype, exDeleter)
	return ex, nil
}

func (ex *Exchange) Close() {
	ex.Closed = true
}

func exchangeNameToType(extype string) (uint8, error) {
	switch extype {
	case "direct":
		return EX_DIRECT, nil
	case "fanout":
		return EX_FANOUT, nil
	case "header":
		return EX_HEADERS, nil
	default:
		return 0, fmt.Errorf("unknown exchange type: %s", extype)
	}
}

func (ex *Exchange) AddBinding(b *binding.Binding, connID int64) error {
	ex.bindLock.Lock()
	defer ex.bindLock.Unlock()

	// Check if binding already exixts
	for _, binding := range ex.bindings {
		if b.Equals(binding) {
			return nil
		}
	}

	ex.bindings = append(ex.bindings, b)
	return nil
}

func (ex *Exchange) RemoveBinding(b *binding.Binding) error {
	ex.bindLock.Lock()
	defer ex.bindLock.Unlock()

	for i, bind := range ex.bindings {
		if b.Equals(bind) {
			ex.bindings = append(ex.bindings[:i], ex.bindings[i+1:]...)
		}
		return nil
	}
	return nil
}

func (ex *Exchange) RemoveQueueBindings(qname string) {
	bindings := make([]*binding.Binding, 0)
	ex.bindLock.Lock()
	defer ex.bindLock.Unlock()

	for _, b := range ex.bindings {
		if b.QueueName != qname {
			bindings = append(bindings, b)
		}
	}
	ex.bindings = bindings
}

func (ex *Exchange) QueuesToPublish(msg *proto.Message) ([]string, *proto.ProtoError) {
	clsID, mtdID := msg.Method.MethodIdentifier()
	queues := make([]string, 0)
	if ex.Name != msg.Method.Exchange {
		return queues, proto.NewSoftError(404, "Exchange name MisMatch", clsID, mtdID)
	}

	switch {
	case ex.ExType == EX_DIRECT:
		for _, b := range ex.bindings {
			if b.CheckDirectMatches(msg.Method) {
				queues = append(queues, b.QueueName)
				return queues, nil
			}
		}
	case ex.ExType == EX_FANOUT:
		for _, b := range ex.bindings {
			if b.CheckDirectMatches(msg.Method) {
				queues = append(queues, b.QueueName)
			}
		}
	// TODO:
	// case ex.ExType == EX_HEADERS
	default:
		panic("Exchange type unknown")
	}

	return queues, nil
}
