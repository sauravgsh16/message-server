package proto

type WireFrame struct {
        FrameType uint8
        Channel   uint16
        Payload   []byte
}

type HeaderFrame struct {
        Class    uint16
        BodySize uint64
}

type Message struct {
        Id       int64
        Header   *HeaderFrame
        Payload  []*WireFrame
        Exchange string
        Method   BasicPublish
}

type BasicPublish struct {

}