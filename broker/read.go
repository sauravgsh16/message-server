package broker

import (
        "io"
        "encoding/binary"
)

func (r *reader) ReadFrame() (frame frame, err error) {
        header := make([]byte, 7)

        if _, err := io.ReadFull(r.r, header); err != nil {
                return nil, err
        }

        _ = uint8(header[0])
        _ = binary.BigEndian.Uint16(header[1:3])
        _ = binary.BigEndian.Uint16(header[3:7])
        return nil, nil
}