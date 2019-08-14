package server

import (
        "github.com/sauravgsh16/secoc-third/qserver/proto"
)

func (ch *Channel) connectionRoute(conn *Connection, mf proto.MethodFrame) *proto.ProtoError {
        switch method := mf.(type) {
        case *proto.ConnectionStartOk:
                return connectionStartOk(conn, method)
        case *proto.ConnectionOpen:
                return ch.connectionOpen(conn, method)
        case *proto.ConnectionCloseOk:
                return ch.connectionOpenOk(conn, method)
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
                Version:   0,
                Mechanism: "PLAIN",
        })
}