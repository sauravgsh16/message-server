package server

import (
	"github.com/sauravgsh16/secoc-third/secoc-final/proto"
)

func (ch *Channel) connectionRoute(conn *Connection, msgf proto.MessageFrame) *proto.Error {
	switch method := msgf.(type) {

	case *proto.ConnectionStartOk:
		return ch.connectionStartOk(conn, method)

	case *proto.ConnectionOpen:
		return ch.connectionOpen(conn, method)

	case *proto.ConnectionClose:
		return ch.connectionClose(conn, method)

	case *proto.ConnectionCloseOk:
		return ch.connectionCloseOk(conn, method)

	default:
		clsID, mtdID := method.MethodIdentifier()
		return proto.NewHardError(540, "unable to route frame", clsID, mtdID)
	}
}

func (ch *Channel) startConnection() *proto.Error {
	ch.conn.status.start = true
	ch.Send(&proto.ConnectionStart{
		Version:    1,
		Mechanisms: "PLAIN",
	})
	return nil
}

func (ch *Channel) connectionOpen(c *Connection, m *proto.ConnectionOpen) *proto.Error {
	// TODO : check if m.Host is accessible.
	// If not, then close connection - break
	c.status.open = true
	ch.Send(&proto.ConnectionOpenOk{Response: "Connected"})
	c.status.openOk = true
	return nil
}

func (ch *Channel) connectionStartOk(c *Connection, m *proto.ConnectionStartOk) *proto.Error {
	c.status.startOk = true

	if m.Mechanism != "PLAIN" {
		c.hardClose()
	}
	return nil
}

func (ch *Channel) connectionClose(c *Connection, m *proto.ConnectionClose) *proto.Error {
	ch.Send(&proto.ConnectionCloseOk{})
	c.hardClose()
	return nil
}

func (ch *Channel) connectionCloseOk(c *Connection, m *proto.ConnectionCloseOk) *proto.Error {
	c.hardClose()
	return nil
}