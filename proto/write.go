package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Writer struct defines field of type io.Writer interface
type Writer struct {
	W io.Writer
}

// WriteFrame calls the write method on the Frame type
func (w Writer) WriteFrame(f Frame) error {
	if err := f.Write(w.W); err != nil {
		return err
	}

	if buf, ok := w.W.(*bufio.Writer); ok {
		if err := buf.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func writeFrame(w io.Writer, fType uint8, channel uint16, payload []byte) error {
	bs := make([]byte, 0, 7+len(payload)+2)
	buf := bytes.NewBuffer(bs)

	if err := WriteOctet(buf, fType); err != nil {
		return err
	}

	if err := WriteShort(buf, channel); err != nil {
		return err
	}

	if err := WriteLongStr(buf, fmt.Sprintf("%s", payload)); err != nil {
		return err
	}

	if err := WriteFrameEnd(buf); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// WriteFrameEnd writes the byte to indicate end of frame
func WriteFrameEnd(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, byte(0xCE))
}

// WriteOctet writes a byte of data
func WriteOctet(w io.Writer, b byte) error {
	return binary.Write(w, binary.BigEndian, b)
}

// WriteShort writes 2 bytes of data
func WriteShort(w io.Writer, i uint16) error {
	return binary.Write(w, binary.BigEndian, i)
}

// WriteLong writes 4 bytes of data
func WriteLong(w io.Writer, i uint32) error {
	return binary.Write(w, binary.BigEndian, i)
}

// WriteLongLong writes 8 bytes of data
func WriteLongLong(w io.Writer, i uint64) error {
	return binary.Write(w, binary.BigEndian, i)
}

// WriteByte one byte of data
func WriteByte(w io.Writer, b byte) error {
	return binary.Write(w, binary.BigEndian, b)
}

// WriteShortStr writes a short sting
func WriteShortStr(w io.Writer, s string) error {
	b := []byte(s)

	length := uint8(len(b))

	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return errors.New("could not write byte: " + err.Error())
	}
	if _, err := w.Write(b[:length]); err != nil {
		return err
	}
	return nil
}

// WriteLongStr writes a long string
func WriteLongStr(w io.Writer, s string) error {
	b := []byte(s)

	length := uint32(len(s))
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}
	if _, err := w.Write(b[:length]); err != nil {
		return err
	}
	return nil
}
