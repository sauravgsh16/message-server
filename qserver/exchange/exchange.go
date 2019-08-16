package exchange

import (
        "fmt"
        "sync"
        "github.com/sauravgsh16/secoc-third/qserver/proto"
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

func NewExchangeFromMethod(m *proto.ExchangeDeclare, exDeleter chan *Exchange) (*Exchange, *proto.ProtoError) {
        extype, err := exchangeNameToType(m.Type)
        if err != nil {
                var classId, methodId = m.MethodIdentifier()
                return nil, proto.NewHardError(503, "Invalid exchange type", classId, methodId)
        }

        ex := NewExchange(m.Exchange, extype, exDeleter)
        return ex, nil
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

func (ex *Exchange) Close() {} // IMPLEMENT