package proto

// ***********************
//    CONNECTION FRAMES
// ***********************

type ConnectionStart struct {
        Version    byte
        Mechanism  string
}

type ConnectionStartOk struct {
        Mechanism  string
        Response   string
}

type ConnectionOpen struct {
        Host string
}

type ConnectionOpenOk struct {
        Response string
}

type ConnectionClose struct {
        ReplyCode uint16
        ReplyText string
        ClassId   uint16
        MethodId  uint16
}

type ConnectionCloseOk struct {}


// ***********************
//      CHANNEL FRAMES
// ***********************

type ChannelOpen struct {
        Reserved string
}

type ChannelOpenOk struct {
        Response string
}

type ChannelFlow struct {
        Active bool
}

type ChannelFlowOk struct {
        Active bool
}

type ChannelClose struct {
        ReplyCode uint16
        ReplyText string
        ClassId   uint16
        MethodId  uint16
}

type ChannelCloseOk struct {}

// ***********************
//     EXCHANGE FRAMES
// ***********************

type ExchangeDeclare struct {
        Exchange string
        Type     string
        NoWait   bool
}

type ExchangeDeclareOk struct {}


type ExchangeDelete struct {
        Exchange string
        IfUnused bool
        NoWait   bool
}

type ExchangeDeleteOk struct {}

type ExchangeBind struct {
        Destination string
        Source      string
        RoutingKey  string
        NoWait      bool
}

type ExchangeBindOk struct {}

type ExchangeUnbind struct {
        Destination string
        Source      string
        RoutingKey  string
        NoWait      bool
}

type ExchangeUnbindOk struct {}

// ***********************
//     EXCHANGE FRAMES
// ***********************

type QueueDeclare struct {
        Queue  string
        NoWait bool
}

type QueueDeclareOk struct {
        Queue       string
        MessageCnt  uint32
        ConsumerCnt uint32
}

type QueueBind struct {
        Queue      string
        Exchange   string
        RoutingKey string
        NoWait     bool
}

type QueueBindOk struct {}

type QueueUnbind struct {
        Queue      string
        Exchange   string
        RoutingKey string
}

type QueueUnbindOk struct {}

type QueueDelete struct {
        Queue     string
	IfUnused  bool
        IfEmpty   bool
        NoWait    bool
}

type QueueDeleteOk struct {
        MessageCnt uint32
}