package client

import (
	"encoding/binary"
	"errors"
	"io"
)

type reader struct {
	r io.Reader
}

func (r *reader) ReadFrame() (frame frame, err error) {
	var scratch [7]byte

	if _, err = io.ReadFull(r.r, scratch[:7]); err != nil {
		return
	}

	typ := uint8(scratch[0])
	channel := binary.BigEndian.Uint16(scratch[1:3])
	size := binary.BigEndian.Uint32(scratch[3:7])

	switch typ {
	case 1:
		if frame, err = r.parseMethodFrame(channel, size); err != nil {
			return
		}
	}

	if _, err = io.ReadFull(r.r, scratch[:1]); err != nil {
		return nil, err
	}

	if scratch[0] != 206 {
		return nil, errors.New("Frame ERRORRRRRR!!!!!")
	}

	return
}

func (r *reader) parseMethodFrame(channel uint16, size uint32) (f frame, err error) {
	mf := &methodFrame{
		ChannelId: channel,
	}

	if err = binary.Read(r.r, binary.BigEndian, &mf.ClassId); err != nil {
		return
	}

	if err = binary.Read(r.r, binary.BigEndian, &mf.MethodId); err != nil {
		return
	}

	switch mf.ClassId {

	case 10: // connection
		switch mf.MethodId {

		case 10: // connection start
			//fmt.Println("NextMethod: class:10 method:10")
			method := &ConnectionStart{}
			if err = method.Read(r.r); err != nil {
				return
			}
			mf.Method = method

		case 11: // connection start-ok
			//fmt.Println("NextMethod: class:10 method:11")
			method := &ConnectionStartOk{}
			if err = method.Read(r.r); err != nil {
				return
			}
			mf.Method = method
		}
	}

	return mf, nil
}

func readLongstr(r io.Reader) (v string, err error) {
	var length uint32
	if err = binary.Read(r, binary.BigEndian, &length); err != nil {
		return
	}

	// slices can't be longer than max int32 value
	if length > (^uint32(0) >> 1) {
		return
	}

	bytes := make([]byte, length)
	if _, err = io.ReadFull(r, bytes); err != nil {
		return
	}
	return string(bytes), nil
}

func readShortstr(r io.Reader) (v string, err error) {
	var length uint8
	if err = binary.Read(r, binary.BigEndian, &length); err != nil {
		return
	}

	bytes := make([]byte, length)
	if _, err = io.ReadFull(r, bytes); err != nil {
		return
	}
	return string(bytes), nil
}
