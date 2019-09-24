package exchange

import (
	"fmt"
	"sync"

	"github.com/sauravgsh16/secoc-third/proto"
	"github.com/sauravgsh16/secoc-third/qserver/binding"
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
	Closed     bool
	deleteChan chan *Exchange
}

func NewExchange(name string, extype uint8, deleteChan chan *Exchange) *Exchange {
	return &Exchange{
		Name:       name,
		ExType:     extype,
		bindings:   make([]*binding.Binding, 0),
		deleteChan: deleteChan,
	}
}

func NewExchangeFromMethod(m *proto.ExchangeDeclare, exDeleter chan *Exchange) (*Exchange, *proto.Error) {
	extype, err := GetExType(m.Type)
	if err != nil {
		clsID, mtdID := m.Identifier()
		return nil, proto.NewHardError(503, err.Error(), clsID, mtdID)
	}

	ex := NewExchange(m.Exchange, extype, exDeleter)
	return ex, nil
}

func (ex *Exchange) Close() {
	ex.Closed = true
}

func GetExType(extype string) (uint8, error) {
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

func (ex *Exchange) QueuesToPublish(msg *proto.Message) ([]string, *proto.Error) {
	clsID, mtdID := msg.Method.Identifier()
	queues := make([]string, 0)
	if ex.Name != msg.Method.(*proto.BasicPublish).Exchange {
		return queues, proto.NewSoftError(404, "Exchange name MisMatch", clsID, mtdID)
	}

	switch {
	case ex.ExType == EX_DIRECT:
		for _, b := range ex.bindings {
			if b.CheckDirectMatches(msg.Method.(*proto.BasicPublish)) {
				queues = append(queues, b.QueueName)
				return queues, nil
			}
		}
	case ex.ExType == EX_FANOUT:
		for _, b := range ex.bindings {
			if b.CheckDirectMatches(msg.Method.(*proto.BasicPublish)) {
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
