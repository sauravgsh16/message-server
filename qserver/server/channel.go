package server

import (
        "bytes"
        "fmt"
        "sync"

        "github.com/sauravgsh16/secoc-third/qserver/proto"
)

// Below struct to have their own files
type Consumer struct{}

// ***** END *********

const (
        CH_INIT = iota
        CH_OPEN
        CH_CLOSING
        CH_CLOSED
)

type Channel struct {
        id             uint16
        server         *Server
        incoming       chan *proto.WireFrame
        outgoing       chan *proto.WireFrame
        conn           *Connection
        consumers      map[string]*Consumer
        sendLock       sync.Mutex
        state          uint8
        currentMessage *proto.Message
        flow           bool
        usedQueueName  string
}

func NewChannel(id uint16, conn *Connection) *Channel {
        return &Channel{
                id:        id,
                server:    conn.server,
                incoming:  make(chan *proto.WireFrame),
                outgoing:  conn.outgoing,
                conn:      conn,
                consumers: make(map[string]*Consumer),
                flow:      true,
        }
}

func (ch *Channel) start() {
        if ch.state == 0 {
                ch.state = CH_OPEN
                go ch.startConnection()
        }

        go func() {
                for {
                        if ch.state == CH_CLOSED {
                                break
                        }
                        var err *proto.ProtoError
                        frame := <- ch.incoming
                        switch {
                        case frame.FrameType == uint8(proto.FrameMethod):
                                fmt.Println("routing method") // LOGS
                                err = ch.routeMethod(frame)
                        case frame.FrameType == uint8(proto.FrameHeader):
                                if ch.state != CH_CLOSING {
                                        fmt.Println("handling header") // LOGS
                                        err = ch.handleHeader(frame)
                                }
                        case frame.FrameType == uint8(proto.FrameBody):
                                if ch.state != CH_CLOSING {
                                        fmt.Println("handling body") // LOGS
                                        err = ch.handleBody(frame)
                                }
                        default:
                                err := proto.NewHardError(500, "Unknown frame type: ", 0, 0)
                        }
                        if err != nil {
                                ch.sendError(err)
                        }
                }
        }()
}

func (ch *Channel) SendMethod(m proto.MethodFrame) {
        buf := bytes.NewBuffer([]byte{})
        m.Write(buf)
        ch.outgoing <- &proto.WireFrame{uint8(proto.FrameMethod), ch.id, buf.Bytes()}
}

func (ch *Channel) sendError(err *proto.ProtoError) {
        if err.Soft {
                fmt.Println("Sending channel error: ", err.Msg)
                ch.state = CH_CLOSING
                ch.SendMethod(&proto.ChannelClose{
                        ReplyCode: err.Code,
                        ReplyText: err.Msg,
                        ClassId:   err.Class,
                        MethodId:  err.Method,
                })
        } else {
                ch.conn.closeConnWithError(err)
        }
}

func (ch *Channel) updateFlow(active bool) {
        if ch.flow == active {
                return
        }
        // Change flow to active
        ch.flow = active
        // Ping Consumers to start work again, if possible
        if ch.flow {
                for _, consumer := range ch.consumers {
                        // *********************
                        // NEED TO IMPLEMENT
                        // *********************
                }
        }
}

func (ch *Channel) shutdown() {
        if ch.state == CH_CLOSED {
                fmt.Printf("channel already closed, shutdown performed on %d\n", ch.id)
                return
        }
        ch.state = CH_CLOSED
        // unregister channel from connection
        ch.conn.removeChannel(ch.id)
        // remove any consumer associated with this channel
        for _, consumer := range ch.consumers {
                // ********************
                // NEED TO IMPLEMENT
                // ********************
        }
}

func (ch *Channel) routeMethod(frame *proto.WireFrame) *proto.ProtoError {
        var methodReader = bytes.NewReader(frame.Payload)

        // ***************** IMPLEMENT BELOW **************
        //                   ReadMethod 
        // ************************************************
        var methodFrame, err = proto.ReadMethod(methodReader)  // TODO - TO BE IMPLEMENTED - ReadMethod
        if err != nil {
                return proto.NewHardError(500, err.Error(), 0, 0)
        }
        var classId, methodId = methodFrame.MethodIdentifier()

        // Check if channel is in initial creation state
        if ch.state == CH_INIT && (classId != 20 || methodId != 10) {
                return proto.NewHardError(503, "Open method call on non-open channel", classId, methodId)
        }
        
        // Route methodFrame based on classId
        switch {
        case classId == 10:
                return ch.connectionRoute(ch.conn, methodFrame)
        case classId == 20:
                return ch.channelRoute(methodFrame)
        case classId == 30:
                return ch.exchangeRoute(methodFrame)
        case classId == 50:
                return ch.queueRoute(methodFrame)
        case classId == 60:
                return ch.basicRoute(methodFrame)
        default:
                return proto.NewHardError(540, "Not Implemented", classId, methodId)
        }
        return nil
}

func (ch *Channel) handleHeader(frame *proto.WireFrame) *proto.ProtoError {

        if ch.currentMessage == nil {
                return proto.NewSoftError(500, "unexpected header frame", 0, 0)
        }

        if ch.currentMessage.Header != nil {
                return proto.NewSoftError(500, "unexpected - already seen header", 0, 0)
        }

        var header = &proto.HeaderFrame{}
        var buf = bytes.NewReader(frame.Payload)
        var err = header.Read(buf)
        if err != nil {
                return proto.NewHardError(500, "Error parsing header frame: " + err.Error(), 0, 0)
        }
        ch.currentMessage.Header = header
        return nil
}

func (ch *Channel) handleBody(frame *proto.WireFrame) *proto.ProtoError {

        if ch.currentMessage == nil {
                return proto.NewSoftError(500, "unexpected header frame", 0, 0)
        }

        if ch.currentMessage.Header == nil {
                return proto.NewSoftError(500, "unexpected body frame - no header yet", 0, 0)
        }

        ch.currentMessage.Payload = append(ch.currentMessage.Payload, frame)

        var size = uint64(0)
        for _, body := range ch.currentMessage.Payload {
                size += uint64(len(body.Payload))
        }
        // Message yet to complete
        if size < ch.currentMessage.Header.BodySize {
                return nil
        }

        // LOGIC TO PUBLISH TO EXCHANGE
        //  
        //       GOES HERE
        //
        // ****************************
        return nil
}