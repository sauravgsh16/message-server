package server

import (
	"github.com/sauravgsh16/secoc-third/proto"
	"github.com/sauravgsh16/secoc-third/qserver/exchange"
)

func (ch *Channel) exchangeRoute(msgf proto.MessageFrame) *proto.Error {
	switch m := msgf.(type) {

	case *proto.ExchangeDeclare:
		return ch.exDeclare(m)

	case *proto.ExchangeDelete:
		return ch.exDelete(m)

	case *proto.ExchangeBind:
		return ch.exBind(m)

	case *proto.ExchangeUnbind:
		return ch.exUnbind(m)

	default:
		clsID, mtdID := msgf.Identifier()
		return proto.NewHardError(540, "Unknown Frame", clsID, mtdID)
	}
}

func (ch *Channel) exDeclare(m *proto.ExchangeDeclare) *proto.Error {

	clsID, mtdID := m.Identifier()

	// Check if exchange is already present in Server
	declared, hasEx := ch.server.getExchange(m.Exchange)
	if hasEx {
		// Check if existing exchange and new exchange have different type
		extype, err := exchange.GetExType(m.Type)
		if err != nil {
			return proto.NewHardError(407, "Unsupported exchange type", clsID, mtdID)
		}
		if declared.ExType != extype {
			return proto.NewHardError(406, "Existing and new exchange have different types", clsID, mtdID)
		}

		if declared.Name == m.Exchange {
			if !m.NoWait {
				ch.Send(&proto.ExchangeDeclareOk{})
			}
		}
		return nil
	}

	// Create new exchange
	ex, pErr := exchange.NewExchangeFromMethod(m, ch.server.exchangeDeleter)
	if pErr != nil {
		return pErr
	}

	err := ch.server.addExchange(ex)
	if err != nil {
		return proto.NewSoftError(500, err.Error(), clsID, mtdID)
	}
	if !m.NoWait {
		ch.Send(&proto.ExchangeDeclareOk{})
	}
	return nil
}

func (ch *Channel) exDelete(m *proto.ExchangeDelete) *proto.Error {
	clsID, mtdID := m.Identifier()
	errCode, err := ch.server.deleteExchange(m)
	if err != nil {
		return proto.NewSoftError(errCode, err.Error(), clsID, mtdID)
	}
	if !m.NoWait {
		ch.Send(&proto.ExchangeDeleteOk{})
	}
	return nil
}

func (ch *Channel) exBind(m *proto.ExchangeBind) *proto.Error {
	cls, mtd := m.Identifier()
	return proto.NewHardError(540, "Not Implemented", cls, mtd)
}

func (ch *Channel) exUnbind(m *proto.ExchangeUnbind) *proto.Error {
	cls, mtd := m.Identifier()
	return proto.NewHardError(540, "Not Implemented", cls, mtd)
}
