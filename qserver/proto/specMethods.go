package proto

import (
        "errors"
        "io"
)

// ** ConnectionStart **

func (f *ConnectionStart) MethodIdentifier() (uint16, uint16) {
        return 10, 10
}

func (f *ConnectionStart) FrameType() byte {
        return 1
}

func (f *ConnectionStart) MethodName() string {
        return "ConnectionStart"
}

func (f *ConnectionStart) Read(r io.Reader) (err error) {

        f.Version, err = ReadOctet(r)
        if err != nil {
                return errors.New("could not read version: " + err.Error())
        }

        f.Mechanism, err = ReadLongStr(r)
        if err != nil {
                return errors.New("counld not read mechanism: " + err.Error())
        }
        return nil
}

// COMPLETE WRITERS
func (f *ConnectionStart) Write(w io.Writer) (err error) {
        return
}


// ** ConnectionStartOk **

func (f *ConnectionStartOk) MethodIdentifier() (uint16, uint16) {
        return 10, 11
}

func (f *ConnectionStartOk) FrameType() byte {
        return 1
}

func (f *ConnectionStartOk) MethodName() string {
        return "ConnectionStartOk"
}

func (f *ConnectionStartOk) Read(r io.Reader) (err error) {
        f.Mechanism, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read mechanism - on ConnStartOk: " + err.Error())
        }

        f.Response, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read response string - on ConnStartOk: " + err.Error())
        }
        return nil
}

func (f *ConnectionStartOk) Write(w io.Writer) (err error) {    // IMPLEMENT IT!!
        return
}

// ** ConnectionOpen **

func (f *ConnectionOpen) MethodIdentifier() (uint16, uint16) {
        return 10, 20
}

func (f *ConnectionOpen) FrameType() byte {
        return 1
}

func (f *ConnectionOpen) MethodName() string {
        return "ConnectionOpen"
}

func (f *ConnectionOpen) Read(r io.Reader) (err error) {
        f.Host, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could to read Host: " + err.Error())
        }
        return nil
}

func (f *ConnectionOpen) Write(writer io.Writer) (err error) {    // IMPLEMENT IT!!
        return
}

// ** ConnectionOpenOk **

func (f *ConnectionOpenOk) MethodIdentifier() (uint16, uint16) {
        return 10, 21
}

func (f *ConnectionOpenOk) FrameType() byte {
        return 1
}

func (f *ConnectionOpenOk) MethodName() string {
        return "ConnectionOpenOk"
}

func (f *ConnectionOpenOk) Read(r io.Reader) (err error) {
        f.Response, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could to read Host: " + err.Error())
        }
        return nil
}

func (f *ConnectionOpenOk) Write(writer io.Writer) (err error) {    // IMPLEMENT IT!!
        return
}

// ** ConnectionClose **

func (f *ConnectionClose) MethodIdentifier() (uint16, uint16) {
        return 10, 30
}

func (f *ConnectionClose) FrameType() byte {
        return 1
}

func (f *ConnectionClose) MethodName() string {
        return "ConnectionClose"
}

func (f *ConnectionClose) Read(r io.Reader) (err error) {
        f.ReplyCode, err = ReadShort(r)
        if err != nil {
                return errors.New("could not read reply code: " + err.Error())
        }

        f.ReplyText, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read reply text: " + err.Error())
        }

        f.ClassId, err = ReadShort(r)
        if err != nil {
                return errors.New("could not read classId: " + err.Error())
        }

        f.MethodId, err = ReadShort(r)
        if err != nil {
                return errors.New("could not read MethodId: " + err.Error())
        }
        return
}

func (f *ConnectionClose) Write(writer io.Writer) (err error) {    // IMPLEMENT IT!!
        return
}


// ** ConnectionCloseOk **

func (f *ConnectionCloseOk) MethodIdentifier() (uint16, uint16) {
        return 10, 31
}

func (f *ConnectionCloseOk) FrameType() byte {
        return 1
}

func (f *ConnectionCloseOk) MethodName() string {
        return "ConnectionCloseOk"
}

func (f *ConnectionCloseOk) Read(r io.Reader) (err error) {
        return
}

func (f *ConnectionCloseOk) Write(writer io.Writer) (err error) {
        return
}