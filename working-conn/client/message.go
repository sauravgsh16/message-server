package client

import (
	"encoding/binary"
	"io"
)

type message interface {
	Read(io.Reader) error
	Write(io.Writer) error
	Id() (uint16, uint16)
}

type frame interface {
	write(io.Writer) error
	channel() uint16
}

type methodFrame struct {
	ChannelId uint16
	ClassId   uint16
	MethodId  uint16
	Method    message
}

func (f *methodFrame) channel() uint16 { return f.ChannelId }

type ConnectionStart struct {
	Version    byte
	Mechanisms string
}

type ConnectionStartOk struct {
	Mechanism string
	Response  string
}

func (msg *ConnectionStart) Id() (uint16, uint16) {
	return 10, 10
}

func (msg *ConnectionStart) Write(w io.Writer) (err error) {

	if err = binary.Write(w, binary.BigEndian, msg.Version); err != nil {
		return
	}

	if err = writeLongstr(w, msg.Mechanisms); err != nil {
		return
	}
	return
}

func (msg *ConnectionStart) Read(r io.Reader) (err error) {

	if err = binary.Read(r, binary.BigEndian, &msg.Version); err != nil {
		return
	}
	if msg.Mechanisms, err = readLongstr(r); err != nil {
		return
	}

	return
}

func (msg *ConnectionStartOk) Id() (uint16, uint16) {
	return 10, 11
}

func (msg *ConnectionStartOk) Write(w io.Writer) (err error) {

	if err = writeShortstr(w, msg.Mechanism); err != nil {
		return
	}

	if err = writeLongstr(w, msg.Response); err != nil {
		return
	}

	return
}

func (msg *ConnectionStartOk) Read(r io.Reader) (err error) {

	if msg.Mechanism, err = readShortstr(r); err != nil {
		return
	}

	if msg.Response, err = readLongstr(r); err != nil {
		return
	}

	return
}
