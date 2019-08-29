package server

import (
	"github.com/sauravgsh16/secoc-third/proto"
)

func (ch *Channel) txRoute(mf proto.MethodFrame) *proto.Error {
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
}

func (ch *Channel) txSelect(m *proto.TxSelect) *proto.Error {
	ch.startTxMode()
	ch.SendMethod(&proto.TxSelectOk{})
	return nil
}

func (ch *Channel) txCommit(m *proto.TxCommit) *proto.Error {
	clsID, mtdID := m.MethodIdentifier()
	if err := ch.commitTx(clsID, mtdID); err != nil {
		return err
	}
	ch.SendMethod(&proto.TxCommitOk{})
	return nil
}

func (ch *Channel) txRollback(m *proto.TxRollback) *proto.Error {
	if err := ch.rollbackTx(); err != nil {
		return err
	}
	ch.SendMethod(&proto.TxRollbackOk{})
	return nil
}
