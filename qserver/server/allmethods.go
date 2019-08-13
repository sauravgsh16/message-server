package server

import (
        "github.com/sauravgsh16/secoc-third/qserver/proto"
)

func (ch *Channel) connectionRoute(conn *Connection, mf proto.MethodFrame) *proto.ProtoError {
        switch method := mf.(type) {
        case *proto.ConnectionOpen:
                return ch.connectionOpen(conn, mf)
        }
        classId, methodId := mf.MethodIdentifier()
        return proto.NewHardError(540, "unable to route frame", classId, methodId)
}

func (ch *Channel) connectionOpen(conn *Connection, mf proto.MethodFrame) *proto.ProtoError {
        conn.status.open = true
        ch.SendMethod(&proto.ConnectionOpen{})
        return nil
}