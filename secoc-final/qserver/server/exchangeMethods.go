package server

import (
	"github.com/sauravgsh16/secoc-third/secoc-final/proto"
	"github.com/sauravgsh16/secoc-third/secoc-final/qserver/exchange"
)

func (ch *Channel) exchangeRoute(msgf proto.MessageFrame) *proto.Error {
	switch method := msgf.(type) {

	case *proto.ExchangeDeclare:
		return ch.exchangeDeclare(method)

	case *proto.ExchangeDelete:
		return ch.exchangeDelete(method)

	case *proto.ExchangeBind:
		return ch.exchangeBind(method)

	case *proto.ExchangeUnbind:
		return ch.exchangeUnbind(method)

	default:
		clsID, mtdID := msgf.MethodIdentifier()
		return proto.NewHardError(540, "Unknown Frame", clsID, mtdID)
	}
}

func (ch *Channel) exchangeDeclare(m *proto.ExchangeDeclare) *proto.Error {

	ex, protoErr := exchange.NewExchangeFromMethod(m, ch.server.exchangeDeleter)
	if protoErr != nil {
		return protoErr
	}

	clsId, mtdId := m.MethodIdentifier()
	// Check if exchange is already present in Server
	existing, hasEx := ch.server.exchanges[ex.Name]
	if hasEx {
		// Check if existing exchange and new exchange have different type
		if existing.ExType != ex.ExType {
			if !m.NoWait {
				return proto.NewHardError(406, "Existing and new exchange have different types", clsId, mtdId)
			}
		}
		ch.Send(&proto.ExchangeDeclareOk{})
		return nil
	}
	err := ch.server.addExchange(ex)
	if err != nil {
		return proto.NewSoftError(500, err.Error(), clsId, mtdId)
	}
	if !m.NoWait {
		ch.Send(&proto.ExchangeDeclareOk{})
	}
	return nil
}

func (ch *Channel) exchangeDelete(m *proto.ExchangeDelete) *proto.Error {
	clsId, mtdId := m.MethodIdentifier()
	errCode, err := ch.server.deleteExchange(m)
	if err != nil {
		return proto.NewSoftError(errCode, err.Error(), clsId, mtdId)
	}
	if !m.NoWait {
		ch.Send(&proto.ExchangeDeleteOk{})
	}
	return nil
}

func (ch *Channel) exchangeBind(m *proto.ExchangeBind) *proto.Error {
	cls, mtd := m.MethodIdentifier()
	return proto.NewHardError(540, "Not Implemented", cls, mtd)
}

func (ch *Channel) exchangeUnbind(m *proto.ExchangeUnbind) *proto.Error {
	cls, mtd := m.MethodIdentifier()
	return proto.NewHardError(540, "Not Implemented", cls, mtd)
}
