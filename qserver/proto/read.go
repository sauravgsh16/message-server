package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func ReadFrame(reader io.Reader) (*WireFrame, error) {

	incoming := make([]byte, 1+2+4)
	if err := binary.Read(reader, binary.LittleEndian, incoming); err != nil {
		return nil, err
	}
	f := &WireFrame{}

	var buf = bytes.NewBuffer(incoming)

	// Get the frametype
	frameType, _ := ReadOctet(buf)
	f.FrameType = frameType

	// Get the channel
	channel, _ := ReadShort(buf)
	f.Channel = channel

	// Get the variable length payload
	var length uint32
	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	var slice = make([]byte, length+1)
	if err := binary.Read(reader, binary.BigEndian, slice); err != nil {
		return nil, errors.New("Bad frame payload " + err.Error())
	}
	f.Payload = slice[0:length]
	return f, nil
}

func ReadOctet(buf io.Reader) (data byte, err error) {
	if err = binary.Read(buf, binary.BigEndian, &data); err != nil {
		return 0, errors.New("Could not read byte: " + err.Error())
	}
	return data, nil
}

func ReadShort(buf io.Reader) (data uint16, err error) {
	if err = binary.Read(buf, binary.BigEndian, &data); err != nil {
		return 0, errors.New("Could not read uint16: " + err.Error())
	}
	return data, nil
}

func ReadLong(buf io.Reader) (data uint32, err error) {
	if err = binary.Read(buf, binary.BigEndian, &data); err != nil {
		return 0, errors.New("Could not read uint32: " + err.Error())
	}
	return data, nil
}

func ReadLongLong(buf io.Reader) (data uint64, err error) {
	if err = binary.Read(buf, binary.BigEndian, &data); err != nil {
		return 0, errors.New("could not read uint64: " + err.Error())
	}
	return data, nil
}

func ReadShortStr(buf io.Reader) (data string, err error) {
	var lenght uint8
	if err = binary.Read(buf, binary.BigEndian, &lenght); err != nil {
		return "", err
	}

	slice := make([]byte, lenght)
	if err = binary.Read(buf, binary.BigEndian, slice); err != nil {
		return "", err
	}
	return string(slice), nil
}

func readLongStr(buf io.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	var slice = make([]byte, length)
	if err := binary.Read(buf, binary.BigEndian, slice); err != nil {
		return nil, err
	}
	return slice, nil
}

func ReadLongStr(buf io.Reader) (string, error) {
	slice, err := readLongStr(buf)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", slice), nil
}

func ReadMethod(r io.Reader) (MethodFrame, error) {
	classID, err := ReadShort(r)
	if err != nil {
		return nil, err
	}

	methodID, err := ReadShort(r)
	if err != nil {
		return nil, err
	}

	switch {
	case classID == 10:
		switch {
		case methodID == 10:
			method := &ConnectionStart{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 11:
			method := &ConnectionStartOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		case methodID == 20:
			method := &ConnectionOpen{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 21:
			method := &ConnectionOpenOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 30:
			method := &ConnectionClose{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 31:
			method := &ConnectionCloseOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		}

	case classID == 20:
		switch {
		case methodID == 10:
			method := &ChannelOpen{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 11:
			method := &ChannelOpenOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		case methodID == 20:
			method := &ChannelFlow{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 21:
			method := &ChannelFlowOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 30:
			method := &ChannelClose{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 31:
			method := &ChannelCloseOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		}

	case classID == 30:
		switch {
		case methodID == 10:
			method := &ExchangeDeclare{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 11:
			method := &ExchangeDeclareOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		case methodID == 20:
			method := &ExchangeDelete{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 21:
			method := &ExchangeDeleteOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 30:
			method := &ExchangeBind{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 31:
			method := &ExchangeBindOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 40:
			method := &ExchangeUnbind{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 41:
			method := &ExchangeUnbindOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		}

	case classID == 40:
		switch {
		case methodID == 10:
			method := &QueueDeclare{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 11:
			method := &QueueDeclareOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		case methodID == 20:
			method := &QueueBind{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 21:
			method := &QueueBindOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 30:
			method := &QueueUnbind{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 31:
			method := &QueueUnbindOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 40:
			method := &QueueDelete{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 41:
			method := &QueueDeleteOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		}

	case classID == 50:
		switch {
		case methodID == 10:
			method := &BasicConsume{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 11:
			method := &BasicConsumeOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		case methodID == 20:
			method := &BasicCancel{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 21:
			method := &BasicCancelOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 30:
			method := &BasicPublish{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 40:
			method := &BasicReturn{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 50:
			method := &BasicDeliver{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 60:
			method := &BasicAck{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 70:
			method := &BasicNack{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		}

	case classID == 60:
		switch {
		case methodID == 10:
			method := &TxSelect{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 11:
			method := &TxSelectOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		case methodID == 20:
			method := &TxCommit{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 21:
			method := &TxCommitOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 30:
			method := &TxRollback{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil

		case methodID == 31:
			method := &TxRollbackOk{}
			err = method.Read(r)
			if err != nil {
				return nil, err
			}
			return method, nil
		}
	}

	return nil, fmt.Errorf("Bad class or method id!. Class id: %d, Method id: %d", classID, methodID)
}
