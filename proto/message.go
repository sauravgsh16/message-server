package proto

type WireFrame struct {
	FrameType uint8
	Channel   uint16
	Payload   []byte
}

type ChannelFrame struct {
	ChannelID uint16
	Method    MethodFrame
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
	MsgSize       uint32
}

type TxMessage struct {
	Msg       *Message
	QueueName string
}

type IndexMessage struct {
	ID            int64
	Refs          int32
	DeliveryCount int32
	Persisted     bool
}
