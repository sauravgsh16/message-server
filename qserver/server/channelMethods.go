package server

import (
	"github.com/sauravgsh16/message-server/proto"
)

func (ch *Channel) channelRoute(msgf proto.MessageFrame) *proto.Error {
	switch method := msgf.(type) {

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

	default:
		clsID, mtdID := msgf.Identifier()
		return proto.NewHardError(540, "Unknown Frame", clsID, mtdID)
	}
}

func (ch *Channel) channelOpen(m *proto.ChannelOpen) *proto.Error {
	if ch.state == chOpen {
		clsID, mtdID := m.Identifier()
		return proto.NewHardError(504, "channel already open", clsID, mtdID)
	}
	ch.Send(&proto.ChannelOpenOk{Response: "200"})
	ch.state = chOpen
	return nil
}

func (ch *Channel) channelFlow(m *proto.ChannelFlow) *proto.Error {
	ch.activateFlow(m.Active)
	ch.Send(&proto.ChannelFlowOk{Active: ch.flow})
	return nil
}

func (ch *Channel) channelFlowOk(m *proto.ChannelFlowOk) *proto.Error {
	cls, mtd := m.Identifier()
	return proto.NewHardError(40, "Not Implemented", cls, mtd)
}

func (ch *Channel) channelClose(m *proto.ChannelClose) *proto.Error {
	ch.Send(&proto.ChannelCloseOk{})
	ch.shutdown()
	return nil
}

func (ch *Channel) channelCloseOk(m *proto.ChannelCloseOk) *proto.Error {
	ch.shutdown()
	return nil
}
