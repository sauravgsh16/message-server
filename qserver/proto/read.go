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

func ReadMethod(io.Reader) (m MethodFrame, err error) {
        return
}