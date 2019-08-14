package proto

// Connection Frames

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


// Channel Frames
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