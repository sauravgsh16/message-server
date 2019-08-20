package proto

type WireFrame struct {
	FrameType uint8
	Channel   uint16
	Payload   []byte
}

type HeaderFrame struct {
	Class    uint16
	BodySize uint64
}

type Message struct {
	ID         int64
	Header     *HeaderFrame
	Payload    []*WireFrame
	Exchange   string
	RoutingKey string
	Method     *BasicPublish
}

type QueueMessage struct {
	ID            int64
	DeliveryCount int32
	msgSize       int32
}

func NewMessage(m *BasicPublish) *Message {
	return &Message{
		ID:         NextCnt(),
		Method:     m,
		Exchange:   m.Exchange,
		RoutingKey: m.RoutingKey,
		Payload:    make([]*WireFrame, 0, 1),
	}
}
