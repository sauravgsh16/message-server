package server

import (
	"github.com/sauravgsh16/message-server/proto"
)

func (ch *Channel) connRoute(conn *Connection, msgf proto.MessageFrame) *proto.Error {
	switch m := msgf.(type) {

	case *proto.ConnectionStartOk:
		return ch.connStartOk(conn, m)

	case *proto.ConnectionOpen:
		return ch.connOpen(conn, m)

	case *proto.ConnectionClose:
		return ch.connClose(conn, m)

	case *proto.ConnectionCloseOk:
		return ch.connCloseOk(conn, m)

	default:
		clsID, mtdID := m.Identifier()
		return proto.NewHardError(540, "unable to route frame", clsID, mtdID)
	}
}

func (ch *Channel) connOpen(c *Connection, m *proto.ConnectionOpen) *proto.Error {
	// TODO : check if m.Host is accessible.
	// If not, then close connection - break
	c.status.open = true
	ch.Send(&proto.ConnectionOpenOk{Response: "Connected"})
	c.status.openOk = true
	return nil
}

func (ch *Channel) connStartOk(c *Connection, m *proto.ConnectionStartOk) *proto.Error {
	c.status.startOk = true

	if m.Mechanism != "PLAIN" {
		c.hardClose()
	}
	return nil
}

func (ch *Channel) connClose(c *Connection, m *proto.ConnectionClose) *proto.Error {
	ch.Send(&proto.ConnectionCloseOk{})
	c.mux.Lock()
	defer c.mux.Unlock()

	c.status.closing = true

	return nil
}

func (ch *Channel) connCloseOk(c *Connection, m *proto.ConnectionCloseOk) *proto.Error {
	c.hardClose()
	return nil
}
