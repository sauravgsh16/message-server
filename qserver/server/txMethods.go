package server

import (
	"github.com/sauravgsh16/message-server/proto"
)

func (ch *Channel) txRoute(msgf proto.MessageFrame) *proto.Error {
	switch m := msgf.(type) {

	case *proto.TxSelect:
		return ch.txSelect(m)

	case *proto.TxCommit:
		return ch.txCommit(m)

	case *proto.TxRollback:
		return ch.txRollback(m)

	default:
		clsID, mtdID := msgf.Identifier()
		return proto.NewHardError(540, "unable to route method frame", clsID, mtdID) // ERROR CODE -- IMPLEMENTATION
	}
}

func (ch *Channel) txSelect(m *proto.TxSelect) *proto.Error {
	ch.startTxMode()
	ch.Send(&proto.TxSelectOk{})
	return nil
}

func (ch *Channel) txCommit(m *proto.TxCommit) *proto.Error {
	clsID, mtdID := m.Identifier()
	if err := ch.commitTx(clsID, mtdID); err != nil {
		return err
	}
	ch.Send(&proto.TxCommitOk{})
	return nil
}

func (ch *Channel) txRollback(m *proto.TxRollback) *proto.Error {
	if err := ch.rollbackTx(); err != nil {
		return err
	}
	ch.Send(&proto.TxRollbackOk{})
	return nil
}
