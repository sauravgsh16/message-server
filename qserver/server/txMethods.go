package server

import (
	"github.com/sauravgsh16/secoc-third/qserver/proto"
)

func (ch *Channel) txRoute(mf proto.MethodFrame) *proto.ProtoError {
	switch m := mf.(type) {
	case *proto.TxSelect:
		return ch.txSelect(m)
	case *proto.TxCommit:
		return ch.txCommit(m)
	case *proto.TxRollback:
		return ch.txRollback(m)
	default:
		clsID, mtdID := m.MethodIdentifier()
		return proto.NewHardError(540, "unable to route method frame", clsID, mtdID) // ERROR CODE -- IMPLEMENTATION
	}
	return nil
}

func (ch *Channel) txSelect(m *proto.TxSelect) *proto.ProtoError {
	ch.startTxMode()
	ch.SendMethod(&proto.TxSelectOk{})
	return nil
}

func (ch *Channel) txCommit(m *proto.TxCommit) *proto.ProtoError {
	if err := ch.commitTx(); err != nil {
		return err
	}
	ch.SendMethod(&proto.TxCommitOk{})
	return nil
}

func (ch *Channel) txRollback(m *proto.TxRollback) *proto.ProtoError {
	if err := ch.rollbackTx(); err != nil {
		return err
	}
	ch.SendMethod(&proto.TxRollbackOk{})
	return nil
}
