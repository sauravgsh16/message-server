package broker

import (
        "io"
        "bufio"
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

func (f *bodyFrame) write(w io.Writer) error {
        return writeFrame(w, f.ChannelId, f.Body)
}

func writeFrame(w io.Writer, channel uint16, payload []byte) error {
        return nil
}