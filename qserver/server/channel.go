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
        id       uint64
        server   *Server
        incoming chan *proto.WireFrame
        outgoing chan *proto.WireFrame
        conn     *Connection
        consumer map[string]*Consumer
        sendLock sync.Mutex
        state    uint8
}

func NewChannel(id uint64, conn *Connection) *Channel {
        return &Channel{
                id:       id,
                server:   conn.server,
                incoming: make(chan *proto.WireFrame),
                outgoing: conn.outgoing,
                conn:     conn,
                consumer: make(map[string]*Consumer),
        }
}

func (ch *Channel) startConnection() {

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

func (ch *Channel) routeMethod(frame *proto.WireFrame) *proto.ProtoError {
        var methodReader = bytes.NewReader(frame.Payload)
        var methodFrame, err = proto.ReadMethod(methodReader)
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