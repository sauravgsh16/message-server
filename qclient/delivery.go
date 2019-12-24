package qclient

import (
	"github.com/sauravgsh16/message-server/proto"
)

// Delivery struct
type Delivery struct {
	ConsumerTag string
	DeliveryTag uint64
	Exchange    string
	RoutingKey  string
	Body        []byte
}

func newDelivery(ch *Channel, mcf proto.MessageContentFrame) *Delivery {
	body := mcf.GetBody()
	d := &Delivery{
		Body: body,
	}

	switch m := mcf.(type) {
	case *proto.BasicDeliver:
		d.ConsumerTag = m.ConsumerTag
		d.DeliveryTag = m.DeliveryTag
		d.Exchange = m.Exchange
		d.RoutingKey = m.RoutingKey
	}
	return d
}
