package broker

import (
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