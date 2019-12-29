package proto

// ***********************
//    CONNECTION FRAMES
// ***********************

// ConnectionStart struct
type ConnectionStart struct {
	Version    byte
	Mechanisms string
}

// ConnectionStartOk struct
type ConnectionStartOk struct {
	Mechanism string
	Response  string
}

// ConnectionOpen struct
type ConnectionOpen struct {
	Host string
}

// ConnectionOpenOk struct
type ConnectionOpenOk struct {
	Response string
}

// ConnectionClose struct
type ConnectionClose struct {
	ReplyCode uint16
	ReplyText string
	ClassId   uint16
	MethodId  uint16
}

// ConnectionCloseOk struct
type ConnectionCloseOk struct{}

// ***********************
//      CHANNEL FRAMES
// ***********************

// ChannelOpen struct
type ChannelOpen struct {
	Reserved string
}

// ChannelOpenOk struct
type ChannelOpenOk struct {
	Response string
}

// ChannelFlow struct
type ChannelFlow struct {
	Active bool
}

// ChannelFlowOk struct
type ChannelFlowOk struct {
	Active bool
}

// ChannelClose struct
type ChannelClose struct {
	ReplyCode uint16
	ReplyText string
	ClassId   uint16
	MethodId  uint16
}

// ChannelCloseOk struct
type ChannelCloseOk struct{}

// ***********************
//     EXCHANGE FRAMES
// ***********************

// ExchangeDeclare struct
type ExchangeDeclare struct {
	Exchange string
	Type     string
	NoWait   bool
}

// ExchangeDeclareOk struct
type ExchangeDeclareOk struct{}

// ExchangeDelete struct
type ExchangeDelete struct {
	Exchange string
	IfUnused bool
	NoWait   bool
}

// ExchangeDeleteOk struct
type ExchangeDeleteOk struct{}

// ExchangeBind struct
type ExchangeBind struct {
	Destination string
	Source      string
	RoutingKey  string
	NoWait      bool
}

// ExchangeBindOk struct
type ExchangeBindOk struct{}

// ExchangeUnbind struct
type ExchangeUnbind struct {
	Destination string
	Source      string
	RoutingKey  string
	NoWait      bool
}

// ExchangeUnbindOk struct
type ExchangeUnbindOk struct{}

// ***********************
//     EXCHANGE FRAMES
// ***********************

// QueueDeclare struct
type QueueDeclare struct {
	Queue  string
	NoWait bool
}

// QueueDeclareOk struct
type QueueDeclareOk struct {
	Queue       string
	MessageCnt  uint32
	ConsumerCnt uint32
}

// QueueBind struct
type QueueBind struct {
	Queue      string
	Exchange   string
	RoutingKey string
	NoWait     bool
}

// QueueBindOk struct
type QueueBindOk struct{}

// QueueUnbind struct
type QueueUnbind struct {
	Queue      string
	Exchange   string
	RoutingKey string
}

// QueueUnbindOk struct
type QueueUnbindOk struct{}

// QueueDelete struct
type QueueDelete struct {
	Queue    string
	IfUnused bool
	IfEmpty  bool
	NoWait   bool
}

// QueueDeleteOk struct
type QueueDeleteOk struct {
	MessageCnt uint32
}

// ***********************
//     BASIC FRAMES
// ***********************

// BasicConsume struct
type BasicConsume struct {
	Queue       string
	ConsumerTag string
	NoAck       bool
	NoWait      bool
}

// BasicConsumeOk struct
type BasicConsumeOk struct {
	ConsumerTag string
}

// BasicCancel struct
type BasicCancel struct {
	ConsumerTag string
	NoWait      bool
}

// BasicCancelOk struct
type BasicCancelOk struct {
	ConsumerTag string
}

// BasicPublish struct
type BasicPublish struct {
	Exchange   string
	RoutingKey string
	Immediate  bool
	Properties Properties
	Body       []byte
}

// BasicReturn struct
type BasicReturn struct {
	ReplyCode  uint16
	ReplyText  string
	Exchange   string
	RoutingKey string
	Properties Properties
	Body       []byte
}

// BasicDeliver struct
type BasicDeliver struct {
	ConsumerTag string
	DeliveryTag uint64
	Exchange    string
	RoutingKey  string
	Properties  Properties
	Body        []byte
}

// BasicAck struct
type BasicAck struct {
	DeliveryTag uint64
	Multiple    bool
}

// BasicNack struct
type BasicNack struct {
	DeliveryTag uint64
	Multiple    bool
	Requeue     bool
}

// ***********************
//    	TX FRAMES
// ***********************

// TxSelect struct
type TxSelect struct{}

// TxSelectOk struct
type TxSelectOk struct{}

// TxCommit struct
type TxCommit struct{}

// TxCommitOk struct
type TxCommitOk struct{}

// TxRollback struct
type TxRollback struct{}

// TxRollbackOk struct
type TxRollbackOk struct{}
