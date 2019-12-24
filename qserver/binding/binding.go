package binding

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/sauravgsh16/message-server/proto"
)

// Binding struct
type Binding struct {
	ID        string
	QueueName string
	Exchange  string
	Key       string
}

// NewBinding returns a binding struct of queue, exchange and routing key
func NewBinding(queue, exchange, key string) (*Binding, error) {
	id, err := calculateID(queue, exchange, key)
	if err != nil {
		return nil, err
	}
	idStr := fmt.Sprintf("%s", id)
	return &Binding{
		QueueName: queue,
		Exchange:  exchange,
		Key:       key,
		ID:        idStr,
	}, nil
}

// Equals check if two bindings are same
func (b *Binding) Equals(b2 *Binding) bool {
	if b.ID != b2.ID {
		return false
	}
	if b.QueueName != b2.QueueName {
		return false
	}
	if b.Exchange != b2.Exchange {
		return false
	}
	if b.Key != b2.Key {
		return false
	}
	return true
}

// CheckDirectMatches checks if the binding exchange and routing key
// matches with the exchange and routing key of the message for direct binding
func (b *Binding) CheckDirectMatches(m *proto.BasicPublish) bool {
	return b.Exchange == m.Exchange && b.Key == m.RoutingKey
}

// CheckFanoutMatches checks if the binding exchange
// matches with the exchange of the message for fanout binding
func (b *Binding) CheckFanoutMatches(m *proto.BasicPublish) bool {
	return b.Exchange == m.Exchange
}

func calculateID(queue, exchange, key string) ([]byte, error) {
	qb := &proto.QueueBind{
		Queue:      queue,
		Exchange:   exchange,
		RoutingKey: key,
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	qb.Write(buf)

	// After writing, we know that the first four bytes will be the classID and MethodID
	val := buf.Bytes()[4:]
	// We then hash the bytes
	hash := sha1.New()
	if _, err := hash.Write(val); err != nil {
		return make([]byte, 0), err
	}
	return []byte(hash.Sum(nil)), nil
}
