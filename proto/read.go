package proto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Reader struct {
	R io.Reader
}

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

func ReadOctet(r io.Reader) (data byte, err error) {
	if err = binary.Read(r, binary.BigEndian, &data); err != nil {
		return 0, errors.New("Could not read byte: " + err.Error())
	}
	return data, nil
}

func ReadShort(r io.Reader) (data uint16, err error) {
	if err = binary.Read(r, binary.BigEndian, &data); err != nil {
		return 0, errors.New("Could not read uint16: " + err.Error())
	}
	return data, nil
}

func ReadLong(r io.Reader) (data uint32, err error) {
	if err = binary.Read(r, binary.BigEndian, &data); err != nil {
		return 0, errors.New("Could not read uint32: " + err.Error())
	}
	return data, nil
}

func ReadLongLong(r io.Reader) (data uint64, err error) {
	if err = binary.Read(r, binary.BigEndian, &data); err != nil {
		return 0, errors.New("could not read uint64: " + err.Error())
	}
	return data, nil
}

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
			method := &ConnectionStart{}
			if err = method.Read(r.R); err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 11:
			method := &ConnectionStartOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		case mf.MethodID == 20:
			method := &ConnectionOpen{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 21:
			method := &ConnectionOpenOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 30:
			method := &ConnectionClose{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 31:
			method := &ConnectionCloseOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		}

	case mf.ClassID == 20:
		switch {
		case mf.MethodID == 10:
			method := &ChannelOpen{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 11:
			method := &ChannelOpenOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		case mf.MethodID == 20:
			method := &ChannelFlow{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 21:
			method := &ChannelFlowOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 30:
			method := &ChannelClose{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 31:
			method := &ChannelCloseOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		}

	case mf.ClassID == 30:
		switch {
		case mf.MethodID == 10:
			method := &ExchangeDeclare{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 11:
			method := &ExchangeDeclareOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		case mf.MethodID == 20:
			method := &ExchangeDelete{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 21:
			method := &ExchangeDeleteOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 30:
			method := &ExchangeBind{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 31:
			method := &ExchangeBindOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 40:
			method := &ExchangeUnbind{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 41:
			method := &ExchangeUnbindOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		}

	case mf.ClassID == 40:
		switch {
		case mf.MethodID == 10:
			method := &QueueDeclare{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 11:
			method := &QueueDeclareOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		case mf.MethodID == 20:
			method := &QueueBind{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 21:
			method := &QueueBindOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 30:
			method := &QueueUnbind{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 31:
			method := &QueueUnbindOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 40:
			method := &QueueDelete{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 41:
			method := &QueueDeleteOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		}

	case mf.ClassID == 50:
		switch {
		case mf.MethodID == 10:
			method := &BasicConsume{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 11:
			method := &BasicConsumeOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		case mf.MethodID == 20:
			method := &BasicCancel{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 21:
			method := &BasicCancelOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 30:
			method := &BasicPublish{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 40:
			method := &BasicReturn{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 50:
			method := &BasicDeliver{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 60:
			method := &BasicAck{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 70:
			method := &BasicNack{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		}

	case mf.ClassID == 60:
		switch {
		case mf.MethodID == 10:
			method := &TxSelect{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 11:
			method := &TxSelectOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		case mf.MethodID == 20:
			method := &TxCommit{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 21:
			method := &TxCommitOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 30:
			method := &TxRollback{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method

		case mf.MethodID == 31:
			method := &TxRollbackOk{}
			err = method.Read(r.R)
			if err != nil {
				return nil, err
			}
			mf.Method = method
		}
	default:
		return nil, fmt.Errorf("Bad class or method id!. Class id: %d, Method id: %d", mf.ClassID, mf.MethodID)

	}

	return mf, nil
}
