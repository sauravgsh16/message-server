package proto

type WireFrame struct {
        FrameType uint8
        Channel   uint16
        Payload   []byte
}