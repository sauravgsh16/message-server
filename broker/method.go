package broker

import (
        "io"
)

type exchangeDeclare struct {
        Exchange string
        Type     string
}

type exchangeDeclareOk struct {}

type basicPublish struct {
        Exchange   string
        RoutingKey string
        Body       []byte
        Properties properties
}

func (msg *basicPublish) read(io.Reader) error {
        // TODO
        return nil
}

func (msg *basicPublish) write(io.Writer) error {
        // TODO
        return nil
}

func (msg *basicPublish) getContent() (properties, []byte) {
        return msg.Properties, msg.Body
}

func (msg *basicPublish) setContent(p properties, b []byte) {
        msg.Properties = p
        msg.Body = b
}

type channelCloseOk struct {}

func (msg *channelCloseOk) read(io.Reader) error {
        // TODO
        return nil
}

func (msg *channelCloseOk) write(io.Writer) error {
        // TODO
        return nil
}