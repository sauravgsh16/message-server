package broker

import (
        "io"
        "bufio"
        "bytes"
        "errors"
)


func (w *writer) WriteFrame(f frame) error {
        if err := f.write(w.w); err != nil {
                return err
        }

        if buf, ok := w.w.(*bufio.Writer); ok {
                if err := buf.Flush(); err != nil {
                        return err
                }
        }
        return nil
}

func (f *methodFrame) write(w io.Writer) error {
        if f.Method == nil {
                return errors.New("frame malformed: missing method")
        }
        var payload bytes.Buffer

        if err := f.Method.write(&payload); err != nil {
                return err
        }
        return nil
}

func (f *headerFrame) write(w io.Writer) error {
        // TODO
        return nil
}

func (f *bodyFrame) write(w io.Writer) error {
        return writeFrame(w, f.ChannelId, f.Body)
}

func writeFrame(w io.Writer, channel uint16, payload []byte) error {
        size := uint(len(payload))
        end := []byte{frameEnd}

        _, err := w.Write([]byte{
                byte((channel & 0xff00) >> 8),
                byte((channel & 0x00ff) >> 0),
                byte((size & 0xff000000) >> 24),
                byte((size & 0x00ff0000) >> 16),
                byte((size & 0x0000ff00) >> 8),
                byte((size & 0x000000ff) >> 0),
        })
        if err != nil {
                return err
        }

        if _, err := w.Write(payload); err != nil {
                return err
        }

        if _, err := w.Write(end); err != nil {
                return err
        }
        return nil
}