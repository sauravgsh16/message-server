package proto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Reader struct
type Reader struct {
	R io.Reader
}

// ReadFrame read the frame from the connection and dispatches call
// according to the Frame type
func (r Reader) ReadFrame() (frame Frame, err error) {

	var incoming [7]byte

	if _, err := io.ReadFull(r.R, incoming[:7]); err != nil {
		return nil, err
	}

	fType := uint8(incoming[0])
	channel := binary.BigEndian.Uint16(incoming[1:3])
	size := binary.BigEndian.Uint32(incoming[3:7])

	switch fType {

	case FrameMethod:
		if frame, err = r.readMethod(channel, size); err != nil {
			return nil, err
		}

	case FrameHeader:
		if frame, err = r.readHeader(channel, size); err != nil {
			return nil, err
		}

	case FrameBody:
		if frame, err = r.readBody(channel, size); err != nil {
			return nil, err
		}

	default:
		return nil, NewHardError(FrameErr, "Frame could not be parsed", 0, 0)
	}

	if _, err := io.ReadFull(r.R, incoming[:1]); err != nil {
		return nil, err
	}

	if incoming[0] != FrameEnd {
		return nil, NewHardError(FrameErr, "Frame could not be parsed", 0, 0)
	}

	return frame, nil
}

func propertySet(mask uint8, property int8) bool {
	return int8(mask)&property > 0
}

func (r Reader) readHeader(channel uint16, size uint32) (Frame, error) {
	hf := &HeaderFrame{
		ChannelID: channel,
	}

	if err := binary.Read(r.R, binary.BigEndian, &hf.Class); err != nil {
		return nil, err
	}

	if err := binary.Read(r.R, binary.BigEndian, &hf.BodySize); err != nil {
		return nil, err
	}

	var flags uint8
	var err error

	if err := binary.Read(r.R, binary.BigEndian, &flags); err != nil {
		return nil, err
	}

	if propertySet(flags, flagContentType) {
		if hf.Properties.ContentType, err = ReadShortStr(r.R); err != nil {
			return nil, err
		}
	}

	if propertySet(flags, flagMessageID) {
		if hf.Properties.MessageID, err = ReadShortStr(r.R); err != nil {
			return nil, err
		}
	}

	if propertySet(flags, flagUserID) {
		if hf.Properties.UserID, err = ReadShortStr(r.R); err != nil {
			return nil, err
		}
	}

	if propertySet(flags, flagAppID) {
		if hf.Properties.ApplicationID, err = ReadShortStr(r.R); err != nil {
			return nil, err
		}
	}

	return hf, nil
}

func (r Reader) readBody(channel uint16, size uint32) (Frame, error) {
	bf := &BodyFrame{
		ChannelID: channel,
		Body:      make([]byte, size),
	}

	if _, err := io.ReadFull(r.R, bf.Body); err != nil {
		return nil, err
	}
	return bf, nil
}

// ReadOctet reads 1 byte of data
func ReadOctet(r io.Reader) (data byte, err error) {
	if err = binary.Read(r, binary.BigEndian, &data); err != nil {
		return 0, errors.New("Could not read byte: " + err.Error())
	}
	return data, nil
}

// ReadShort reads 2 bytes of data
func ReadShort(r io.Reader) (data uint16, err error) {
	if err = binary.Read(r, binary.BigEndian, &data); err != nil {
		return 0, errors.New("Could not read uint16: " + err.Error())
	}
	return data, nil
}

// ReadLong reads 4 bytes of data
func ReadLong(r io.Reader) (data uint32, err error) {
	if err = binary.Read(r, binary.BigEndian, &data); err != nil {
		return 0, errors.New("Could not read uint32: " + err.Error())
	}
	return data, nil
}

// ReadLongLong reads a long string
func ReadLongLong(r io.Reader) (data uint64, err error) {
	if err = binary.Read(r, binary.BigEndian, &data); err != nil {
		return 0, errors.New("could not read uint64: " + err.Error())
	}
	return data, nil
}

// ReadShortStr reads a stort string
func ReadShortStr(r io.Reader) (data string, err error) {
	var lenght uint8
	if err = binary.Read(r, binary.BigEndian, &lenght); err != nil {
		return "", err
	}

	slice := make([]byte, lenght)
	if _, err = io.ReadFull(r, slice); err != nil {
		return "", err
	}
	return string(slice), nil
}

func readLongStr(r io.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	var slice = make([]byte, length)
	if _, err := io.ReadFull(r, slice); err != nil {
		return nil, err
	}
	return slice, nil
}

// ReadLongStr reads a string of 4 bytes
func ReadLongStr(r io.Reader) (string, error) {
	slice, err := readLongStr(r)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", slice), nil
}

func (r Reader) readMethod(channelID uint16, size uint32) (Frame, error) {
	mf := &MethodFrame{
		ChannelID: channelID,
	}

	var err error
	var method MessageFrame

	if mf.ClassID, err = ReadShort(r.R); err != nil {
		return nil, err
	}

	if mf.MethodID, err = ReadShort(r.R); err != nil {
		return nil, err
	}

	switch {
	case mf.ClassID == 10:
		switch {
		case mf.MethodID == 10:
			method = &ConnectionStart{}

		case mf.MethodID == 11:
			method = &ConnectionStartOk{}

		case mf.MethodID == 20:
			method = &ConnectionOpen{}

		case mf.MethodID == 21:
			method = &ConnectionOpenOk{}

		case mf.MethodID == 30:
			method = &ConnectionClose{}

		case mf.MethodID == 31:
			method = &ConnectionCloseOk{}
		}

	case mf.ClassID == 20:
		switch {
		case mf.MethodID == 10:
			method = &ChannelOpen{}

		case mf.MethodID == 11:
			method = &ChannelOpenOk{}

		case mf.MethodID == 20:
			method = &ChannelFlow{}

		case mf.MethodID == 21:
			method = &ChannelFlowOk{}

		case mf.MethodID == 30:
			method = &ChannelClose{}

		case mf.MethodID == 31:
			method = &ChannelCloseOk{}
		}

	case mf.ClassID == 30:
		switch {
		case mf.MethodID == 10:
			method = &ExchangeDeclare{}

		case mf.MethodID == 11:
			method = &ExchangeDeclareOk{}

		case mf.MethodID == 20:
			method = &ExchangeDelete{}

		case mf.MethodID == 21:
			method = &ExchangeDeleteOk{}

		case mf.MethodID == 30:
			method = &ExchangeBind{}

		case mf.MethodID == 31:
			method = &ExchangeBindOk{}

		case mf.MethodID == 40:
			method = &ExchangeUnbind{}

		case mf.MethodID == 41:
			method = &ExchangeUnbindOk{}
		}

	case mf.ClassID == 40:
		switch {
		case mf.MethodID == 10:
			method = &QueueDeclare{}

		case mf.MethodID == 11:
			method = &QueueDeclareOk{}

		case mf.MethodID == 20:
			method = &QueueBind{}

		case mf.MethodID == 21:
			method = &QueueBindOk{}

		case mf.MethodID == 30:
			method = &QueueUnbind{}

		case mf.MethodID == 31:
			method = &QueueUnbindOk{}

		case mf.MethodID == 40:
			method = &QueueDelete{}

		case mf.MethodID == 41:
			method = &QueueDeleteOk{}
		}

	case mf.ClassID == 50:
		switch {
		case mf.MethodID == 10:
			method = &BasicConsume{}

		case mf.MethodID == 11:
			method = &BasicConsumeOk{}
		case mf.MethodID == 20:
			method = &BasicCancel{}

		case mf.MethodID == 21:
			method = &BasicCancelOk{}

		case mf.MethodID == 30:
			method = &BasicPublish{}

		case mf.MethodID == 40:
			method = &BasicReturn{}

		case mf.MethodID == 50:
			method = &BasicDeliver{}

		case mf.MethodID == 60:
			method = &BasicAck{}

		case mf.MethodID == 70:
			method = &BasicNack{}
		}

	case mf.ClassID == 60:
		switch {
		case mf.MethodID == 10:
			method = &TxSelect{}

		case mf.MethodID == 11:
			method = &TxSelectOk{}

		case mf.MethodID == 20:
			method = &TxCommit{}

		case mf.MethodID == 21:
			method = &TxCommitOk{}

		case mf.MethodID == 30:
			method = &TxRollback{}

		case mf.MethodID == 31:
			method = &TxRollbackOk{}
		}
	default:
		return nil, fmt.Errorf("Bad class or method id!. Class id: %d, Method id: %d", mf.ClassID, mf.MethodID)

	}

	err = method.Read(r.R)
	if err != nil {
		return nil, err
	}
	mf.Method = method

	return mf, nil
}
