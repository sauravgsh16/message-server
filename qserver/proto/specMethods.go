package proto

import (
        "errors"
        "io"
)

// *********************
//    CONNECTION SPECS
// *********************

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

func (f *ConnectionClose) Write(writer io.Writer) (err error) { // IMPLEMENT IT!!
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


// *******************
//    CHANNEL SPECS
// *******************

// ChannelOpen

func (f *ChannelOpen) MethodIdentifier() (uint16, uint16) {
        return 20, 10
}

func (f *ChannelOpen) MethodName() string {
        return "ChannelOpen"
}

func (f *ChannelOpen) FrameType() byte {
        return 1
}

func (f *ChannelOpen) Read(r io.Reader) (err error) {
        f.Reserved, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read reserved string from channel open: " + err.Error())
        }
        return
}

func (f *ChannelOpen) Write(w io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ChannelOpenOk

func (f *ChannelOpenOk) MethodIdentifier() (uint16, uint16) {
        return 20, 11
}

func (f *ChannelOpenOk) MethodName() string {
        return "ChannelOpenOk"
}

func (f *ChannelOpenOk) FrameType() byte {
        return 1
}

func (f *ChannelOpenOk) Read(r io.Reader) (err error) {
        f.Response, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read response of channel open ok: " + err.Error())
        }
        return
}

func (f *ChannelOpenOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ChannelFlow

func (f *ChannelFlow) MethodIdentifier() (uint16, uint16) {
        return 20, 20
}

func (f *ChannelFlow) MethodName() string {
        return "ChannelFlow"
}

func (f *ChannelFlow) FrameType() byte {
        return 1
}

func (f *ChannelFlow) Read(r io.Reader) (err error) {
        bits, err := ReadOctet(r)
        if err != nil {
                return errors.New("could not read response of channel open ok: " + err.Error())
        }
        f.Active = (bits&(1<<0) > 0)
        return
}

func (f *ChannelFlow) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ChannelFlowOk

func (f *ChannelFlowOk) MethodIdentifier() (uint16, uint16) {
        return 20, 21
}

func (f *ChannelFlowOk) MethodName() string {
        return "ChannelFlowOk"
}

func (f *ChannelFlowOk) FrameType() byte {
        return 1
}

func (f *ChannelFlowOk) Read(r io.Reader) (err error) {
        bits, err := ReadOctet(r)
        if err != nil {
                return errors.New("could not read response of channel open ok: " + err.Error())
        }
        f.Active = (bits&(1<<0) > 0)
        return
}

func (f *ChannelFlowOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ChannelClose

func (f *ChannelClose) MethodIdentifier() (uint16, uint16) {
        return 20, 30
}

func (f *ChannelClose) MethodName() string {
        return "ChannelClose"
}

func (f *ChannelClose) FrameType() byte {
        return 1
}

func (f *ChannelClose) Read(r io.Reader) (err error) {
        f.ReplyCode, err = ReadShort(r)
        if err != nil {
                return errors.New("could not read reply code of channel close: " + err.Error())
        }

        f.ReplyText, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read reply text of channel close: " + err.Error())
        }

        f.ClassId, err = ReadShort(r)
        if err != nil {
                return errors.New("could not read class id of channel close: " + err.Error())
        }

        f.MethodId, err = ReadShort(r)
        if err != nil {
                return errors.New("could not read method id of channel close: " + err.Error())
        }
        return
}

func (f *ChannelClose) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ChannelCloseOk

func (f *ChannelCloseOk) MethodIdentifier() (uint16, uint16) {
        return 20, 31
}

func (f *ChannelCloseOk) MethodName() string {
        return "ChannelCloseOk"
}

func (f *ChannelCloseOk) FrameType() byte {
        return 1
}

func (f *ChannelCloseOk) Read(r io.Reader) (err error) {
        return
}

func (f *ChannelCloseOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}


// *******************
//    EXCHANGE SPECS
// *******************

// ExchangeDeclare

func (f *ExchangeDeclare) MethodIdentifier() (uint16, uint16) {
        return 30, 10
}

func (f *ExchangeDeclare) MethodName() string {
        return "ExchangeDeclare"
}

func (f *ExchangeDeclare) FrameType() byte {
        return 1
}

func (f *ExchangeDeclare) Read(r io.Reader) (err error) {
        f.Exchange, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read exchange name in declare: " + err.Error())
        }

        f.Type, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read exchange type in declare: " + err.Error())
        }

        bits, err := ReadOctet(r)
        if err != nil {
                return errors.New("could not read bits in declare: " + err.Error())
        }
        f.NoWait = (bits&(1<<0) > 0)

        return
}

func (f *ExchangeDeclare) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ExchangeDeclareOk

func (f *ExchangeDeclareOk) MethodIdentifier() (uint16, uint16) {
        return 30, 11
}

func (f *ExchangeDeclareOk) MethodName() string {
        return "ExchangeDeclareOk"
}

func (f *ExchangeDeclareOk) FrameType() byte {
        return 1
}

func (f *ExchangeDeclareOk) Read(r io.Reader) (err error) {
        return
}

func (f *ExchangeDeclareOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ExchangeDelete

func (f *ExchangeDelete) MethodIdentifier() (uint16, uint16) {
        return 30, 20
}

func (f *ExchangeDelete) MethodName() string {
        return "ExchangeDelete"
}

func (f *ExchangeDelete) FrameType() byte {
        return 1
}

func (f *ExchangeDelete) Read(r io.Reader) (err error) {
        f.Exchange, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read exchange name in delete: " + err.Error())
        }

        bits, err := ReadOctet(r)
        if err != nil {
                return errors.New("could not read bits in delete: " + err.Error())
        }

        f.IfUnused = (bits&(1<<0) > 0)
        f.NoWait = (bits&(1<<1) > 0)

        return
}

func (f *ExchangeDelete) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ExchangeDeleteOk

func (f *ExchangeDeleteOk) MethodIdentifier() (uint16, uint16) {
        return 30, 21
}

func (f *ExchangeDeleteOk) MethodName() string {
        return "ExchangeDeleteOk"
}

func (f *ExchangeDeleteOk) FrameType() byte {
        return 1
}

func (f *ExchangeDeleteOk) Read(r io.Reader) (err error) {
        return
}

func (f *ExchangeDeleteOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ExchangeBind

func (f *ExchangeBind) MethodIdentifier() (uint16, uint16) {
        return 30, 30
}

func (f *ExchangeBind) MethodName() string {
        return "ExchangeBind"
}

func (f *ExchangeBind) FrameType() byte {
        return 1
}

func (f *ExchangeBind) Read(r io.Reader) (err error) {
        f.Destination, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read Destination in bind: " + err.Error())
        }

        f.Source, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read Source in bind: " + err.Error())
        }

        f.RoutingKey , err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read RoutingKey in bind: " + err.Error())
        }

        bits, err := ReadOctet(r)
        if err != nil {
                return errors.New("could not read bits in bind: " + err.Error())
        }

        f.NoWait = (bits&(1<<0) > 0)

        return
}

func (f *ExchangeBind) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ExchangeBindOk

func (f *ExchangeBindOk) MethodIdentifier() (uint16, uint16) {
        return 30, 31
}

func (f *ExchangeBindOk) MethodName() string {
        return "ExchangeBindOk"
}

func (f *ExchangeBindOk) FrameType() byte {
        return 1
}

func (f *ExchangeBindOk) Read(r io.Reader) (err error) {
        return
}

func (f *ExchangeBindOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ExchangeUnbind

func (f *ExchangeUnbind) MethodIdentifier() (uint16, uint16) {
        return 30, 40
}

func (f *ExchangeUnbind) MethodName() string {
        return "ExchangeUnbind"
}

func (f *ExchangeUnbind) FrameType() byte {
        return 1
}

func (f *ExchangeUnbind) Read(r io.Reader) (err error) {
        f.Destination, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read Destination in bind: " + err.Error())
        }

        f.Source, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read Source in bind: " + err.Error())
        }

        f.RoutingKey , err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read RoutingKey in bind: " + err.Error())
        }

        bits, err := ReadOctet(r)
        if err != nil {
                return errors.New("could not read bits in bind: " + err.Error())
        }

        f.NoWait = (bits&(1<<0) > 0)

        return
}

func (f *ExchangeUnbind) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// ExchangeUnbindOk

func (f *ExchangeUnbindOk) MethodIdentifier() (uint16, uint16) {
        return 30, 41
}

func (f *ExchangeUnbindOk) MethodName() string {
        return "ExchangeUnbindOk"
}

func (f *ExchangeUnbindOk) FrameType() byte {
        return 1
}

func (f *ExchangeUnbindOk) Read(r io.Reader) (err error) {
        return
}

func (f *ExchangeUnbindOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}


// *******************
//    QUEUE SPECS
//        QueueDeclare - 10
//        QueueDeclareOk - 11
//        QueueBind      - 20
//        QueueBindOk    - 21
//        QueueUnbind    - 30
//        QueueUnbindOk  - 31
//        QueueDelete    - 40
//        QueueDeleteOk  - 41
// *******************

// QueueDeclare

func (f *QueueDeclare) MethodIdentifier() (uint16, uint16) {
        return 40, 10
}


func (f *QueueDeclare) MethodName() string {
        return "QueueDeclare"
}

func (f *QueueDeclare) FrameType() byte {
        return 1
}

func (f *QueueDeclare) Read(r io.Reader) (err error) {
        f.Queue, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read queue name in queue declare: " + err.Error())
        }

        bits, err := ReadOctet(r)
        if err != nil {
                return errors.New("could not read bits in queue declare: " + err.Error())
        }
        f.NoWait = (bits&(1<<0) > 0)

        return
}

func (f *QueueDeclare) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// QueueDeclareOk

func (f *QueueDeclareOk) MethodIdentifier() (uint16, uint16) {
        return 40, 11
}


func (f *QueueDeclareOk) MethodName() string {
        return "QueueDeclareOk"
}

func (f *QueueDeclareOk) FrameType() byte {
        return 1
}

func (f *QueueDeclareOk) Read(r io.Reader) (err error) {
        f.Queue, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read queue name in queue declareOk: " + err.Error())
        }

        f.MessageCnt, err = ReadLong(r)
        if err != nil {
                return errors.New("could not read message count in queue declareOk: " + err.Error())
        }

        f.ConsumerCnt, err = ReadLong(r)
        if err != nil {
                return errors.New("could not read message count in queue declareOk: " + err.Error())
        }

        return
}

func (f *QueueDeclareOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// QueueBind

func (f *QueueBind) MethodIdentifier() (uint16, uint16) {
        return 40, 20
}


func (f *QueueBind) MethodName() string {
        return "QueueBind"
}

func (f *QueueBind) FrameType() byte {
        return 1
}

func (f *QueueBind) Read(r io.Reader) (err error) {
        f.Queue, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read queue name in queue bind: " + err.Error())
        }

        f.Exchange, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read exchange in queue bind: " + err.Error())
        }

        f.RoutingKey, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read routingkey in queue bind: " + err.Error())
        }
        
        bits, err := ReadOctet(r)
        if err != nil {
                return errors.New("could not read bits in queue bind: " + err.Error())
        }
        f.NoWait = (bits&(1<<0) > 0)

        return
}

func (f *QueueBind) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// QueueBindOk

func (f *QueueBindOk) MethodIdentifier() (uint16, uint16) {
        return 40, 21
}


func (f *QueueBindOk) MethodName() string {
        return "QueueBindOk"
}

func (f *QueueBindOk) FrameType() byte {
        return 1
}

func (f *QueueBindOk) Read(r io.Reader) (err error) {
        return
}

func (f *QueueBindOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// QueueUnbind

func (f *QueueUnbind) MethodIdentifier() (uint16, uint16) {
        return 40, 30
}


func (f *QueueUnbind) MethodName() string {
        return "QueueUnbind"
}

func (f *QueueUnbind) FrameType() byte {
        return 1
}

func (f *QueueUnbind) Read(r io.Reader) (err error) {
        f.Queue, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read queue name in queue unbind: " + err.Error())
        }

        f.Exchange, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read exchange in queue unbind: " + err.Error())
        }

        f.RoutingKey, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read routingkey in queue unbind: " + err.Error())
        }

        return
}

func (f *QueueUnbind) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// QueueUnbindOk

func (f *QueueUnbindOk) MethodIdentifier() (uint16, uint16) {
        return 40, 31
}

func (f *QueueUnbindOk) MethodName() string {
        return "QueueUnbindOk"
}

func (f *QueueUnbindOk) FrameType() byte {
        return 1
}

func (f *QueueUnbindOk) Read(r io.Reader) (err error) {
        return
}

func (f *QueueUnbindOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// QueueDelete

func (f *QueueDelete) MethodIdentifier() (uint16, uint16) {
        return 40, 40
}


func (f *QueueDelete) MethodName() string {
        return "QueueDelete"
}

func (f *QueueDelete) FrameType() byte {
        return 1
}

func (f *QueueDelete) Read(r io.Reader) (err error) {
        f.Queue, err = ReadLongStr(r)
        if err != nil {
                return errors.New("could not read queue name in queue delete: " + err.Error())
        }
        
        bits, err := ReadOctet(r)
        if err != nil {
                return errors.New("could not read bits in queue delete: " + err.Error())
        }
        f.IfUnused = (bits&(1<<0) > 0)
        f.IfEmpty = (bits&(1<<1) > 0)
        f.NoWait = (bits&(1<<2) > 0)

        return
}

func (f *QueueDelete) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}

// QueueDeleteOk

func (f *QueueDeleteOk) MethodIdentifier() (uint16, uint16) {
        return 40, 40
}


func (f *QueueDeleteOk) MethodName() string {
        return "QueueDeleteOk"
}

func (f *QueueDeleteOk) FrameType() byte {
        return 1
}

func (f *QueueDeleteOk) Read(r io.Reader) (err error) {
        f.MessageCnt, err = ReadLong(r)
        if err != nil {
                return errors.New("could not read queue name in queue deleteOk: " + err.Error())
        }

        return
}

func (f *QueueDeleteOk) Write(r io.Writer) (err error) { // IMPLEMENT IT!!
        return
}