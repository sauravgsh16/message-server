package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type writer struct {
	w io.Writer
}

func (w *writer) WriteFrame(frame frame) (err error) {
	if err = frame.write(w.w); err != nil {
		return
	}

	if buf, ok := w.w.(*bufio.Writer); ok {
		err = buf.Flush()
	}

	return
}

func (f *methodFrame) write(w io.Writer) (err error) {
	var payload bytes.Buffer

	if f.Method == nil {
		return errors.New("malformed frame: missing method")
	}

	class, method := f.Method.Id()

	if err = binary.Write(&payload, binary.BigEndian, class); err != nil {
		return
	}

	if err = binary.Write(&payload, binary.BigEndian, method); err != nil {
		return
	}

	if err = f.Method.Write(&payload); err != nil {
		return
	}

	return writeFrame(w, 1, f.ChannelId, payload.Bytes())
}

func writeFrame(w io.Writer, typ uint8, channel uint16, payload []byte) (err error) {
	end := []byte{206}
	size := uint(len(payload))

	_, err = w.Write([]byte{
		byte(typ),
		byte((channel & 0xff00) >> 8),
		byte((channel & 0x00ff) >> 0),
		byte((size & 0xff000000) >> 24),
		byte((size & 0x00ff0000) >> 16),
		byte((size & 0x0000ff00) >> 8),
		byte((size & 0x000000ff) >> 0),
	})

	if err != nil {
		return
	}

	if _, err = w.Write(payload); err != nil {
		return
	}

	if _, err = w.Write(end); err != nil {
		return
	}

	return
}

func writeShortstr(w io.Writer, s string) (err error) {
	b := []byte(s)

	var length = uint8(len(b))

	if err = binary.Write(w, binary.BigEndian, length); err != nil {
		return
	}

	if _, err = w.Write(b[:length]); err != nil {
		return
	}

	return
}

func writeLongstr(w io.Writer, s string) (err error) {
	b := []byte(s)

	var length = uint32(len(b))

	if err = binary.Write(w, binary.BigEndian, length); err != nil {
		return
	}

	if _, err = w.Write(b[:length]); err != nil {
		return
	}

	return
}
