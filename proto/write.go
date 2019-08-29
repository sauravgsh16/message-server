package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func WriteFrame(w io.Writer, wf *WireFrame) {
	bs := make([]byte, 0, 7+len(wf.Payload)+2)
	buf := bytes.NewBuffer(bs)

	WriteOctet(buf, wf.FrameType)
	WriteShort(buf, wf.Channel)
	WriteLongStr(buf, fmt.Sprintf("%s", wf.Payload))

	WriteFrameEnd(buf)

	binary.Write(w, binary.LittleEndian, buf.Bytes())
}

func WriteFrameEnd(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, byte(0xCE))
}

func WriteOctet(w io.Writer, b byte) error {
	return binary.Write(w, binary.BigEndian, b)
}

func WriteShort(w io.Writer, i uint16) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteLong(w io.Writer, i uint32) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteLongLong(w io.Writer, i uint64) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteByte(w io.Writer, b byte) error {
	return binary.Write(w, binary.BigEndian, b)
}

func WriteShortStr(w io.Writer, s string) error {
	if len(s) > int(uint8(255)) {
		return errors.New("expected short string: String too long")
	}

	if err := binary.Write(w, binary.BigEndian, byte(len(s))); err != nil {
		return errors.New("could not write byte: " + err.Error())
	}
	return binary.Write(w, binary.BigEndian, s)
}

func WriteLongStr(w io.Writer, s string) error {
	if err := binary.Write(w, binary.BigEndian, uint32(len(s))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, s); err != nil {
		return err
	}
	return nil
}

func WriteMethodIdentifier(w io.Writer, mf MethodFrame) error {
	classID, methodID := mf.MethodIdentifier()

	if err := binary.Write(w, binary.BigEndian, classID); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, methodID); err != nil {
		return err
	}
	return nil
}
