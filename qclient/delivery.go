package qclient

import (
	"github.com/sauravgsh16/message-server/proto"
)

// Delivery struct
type Delivery struct {
	// Properties
	ContentType   string
	MessageID     string
	UserID        string
	ApplicationID string

	ConsumerTag string
	DeliveryTag uint64
	Exchange    string
	RoutingKey  string

	// Payload
	Body []byte
}

func newDelivery(ch *Channel, mcf proto.MessageContentFrame) *Delivery {
	props, body := mcf.GetContent()
	d := &Delivery{
		ContentType:   props.ContentType,
		MessageID:     props.MessageID,
		UserID:        props.UserID,
		ApplicationID: props.ApplicationID,
		Body:          body,
	}

	switch m := mcf.(type) {
	// TODO: switch for BasicGet
	case *proto.BasicDeliver:
		d.ConsumerTag = m.ConsumerTag
		d.DeliveryTag = m.DeliveryTag
		d.Exchange = m.Exchange
		d.RoutingKey = m.RoutingKey
	}
	return d
}
