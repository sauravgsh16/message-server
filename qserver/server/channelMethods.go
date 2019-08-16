package server

import (
        "github.com/sauravgsh16/secoc-third/qserver/proto"
)

func (ch *Channel) channelRoute(mf proto.MethodFrame) *proto.ProtoError {
        switch method := mf.(type) {
        case *proto.ChannelOpen:
                return ch.channelOpen(method)
        case *proto.ChannelFlow:
                return ch.channelFlow(method)
        case *proto.ChannelFlowOk:
                return ch.channelFlowOk(method)
        case *proto.ChannelClose:
                return ch.channelClose(method)
        case *proto.ChannelCloseOk:
                return ch.channelCloseOk(method)
        }
        return nil
}

func (ch *Channel) channelOpen(m *proto.ChannelOpen) *proto.ProtoError {
        if ch.state == CH_OPEN {
                var classId, methodId = m.MethodIdentifier()
                return proto.NewHardError(504, "channel already open", classId, methodId)
        }
        ch.SendMethod(&proto.ChannelOpenOk{})
        ch.state = CH_OPEN
        return nil
}

func (ch *Channel) channelFlow(m *proto.ChannelFlow) *proto.ProtoError {
        ch.updateFlow(m.Active)
        ch.SendMethod(&proto.ChannelFlowOk{ch.flow})
        return nil
}

func (ch *Channel) channelFlowOk(m *proto.ChannelFlowOk) *proto.ProtoError {
        cls, mtd := m.MethodIdentifier()
        return proto.NewHardError(40, "Not Implemented", cls, mtd)
}

func (ch *Channel) channelClose(m *proto.ChannelClose) *proto.ProtoError {
        ch.SendMethod(&proto.ChannelCloseOk{})
        ch.shutdown()
        return nil
}

func (ch *Channel) channelCloseOk(m *proto.ChannelCloseOk) *proto.ProtoError {
        ch.shutdown()
        return nil
}