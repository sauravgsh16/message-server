package server

import (
        "github.com/sauravgsh16/secoc-third/qserver/proto"
)

func (ch *Channel) connectionRoute(conn *Connection, mf proto.MethodFrame) *proto.ProtoError {
        switch method := mf.(type) {
        case *proto.ConnectionStartOk:
                return ch.connectionStartOk(conn, method)
        case *proto.ConnectionOpen:
                return ch.connectionOpen(conn, method)
        case *proto.ConnectionClose:
                return ch.connectionClose(conn, method)
        case *proto.ConnectionCloseOk:
                return ch.connectionCloseOk(conn, method)
        }
        classId, methodId := mf.MethodIdentifier()
        return proto.NewHardError(540, "unable to route frame", classId, methodId)
}

func (ch *Channel) connectionOpen(c *Connection, m *proto.ConnectionOpen) *proto.ProtoError {
        c.status.open = true
        ch.SendMethod(&proto.ConnectionOpenOk{""})
        c.status.openOk = true
        return nil
}

func (ch *Channel) connectionStartOk(c *Connection, m *proto.ConnectionStartOk) *proto.ProtoError {
        c.status.startOk = true

        if m.Mechanism != "PLAIN" {
                c.hardClose()
        }

        return nil
}

func (ch *Channel) startConnection() *proto.ProtoError {
        ch.SendMethod(&proto.ConnectionStart{
                Version:   1,
                Mechanism: "PLAIN",
        })
        return nil
}

func (ch *Channel) connectionClose(c *Connection, m *proto.ConnectionClose) *proto.ProtoError {
        ch.SendMethod(&proto.ConnectionCloseOk{})
        c.hardClose()
        return nil
}

func (ch *Channel) connectionCloseOk(c *Connection, m *proto.ConnectionCloseOk) *proto.ProtoError {
        c.hardClose()
        return nil
}