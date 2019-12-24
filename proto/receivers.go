package proto

import (
	"errors"
	"io"
)

// *********************
//    CONNECTION SPECS
// *********************

// ** ConnectionStart **

// Identifier returns the class ID and method ID
func (f *ConnectionStart) Identifier() (uint16, uint16) {
	return 10, 10
}

// MethodName returns a the name of the Method
func (f *ConnectionStart) MethodName() string {
	return "ConnectionStart"
}

// FrameType returns the frame type of the method
func (f *ConnectionStart) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ConnectionStart) Wait() bool {
	return true
}

func (f *ConnectionStart) Read(r io.Reader) (err error) {

	f.Version, err = ReadOctet(r)
	if err != nil {
		return errors.New("could not read version: " + err.Error())
	}

	f.Mechanisms, err = ReadLongStr(r)
	if err != nil {
		return errors.New("counld not read mechanism: " + err.Error())
	}
	return nil
}

func (f *ConnectionStart) Write(w io.Writer) (err error) {

	if err = WriteByte(w, f.Version); err != nil {
		return errors.New("could not write version in ConnectionStart: " + err.Error())
	}

	if err = WriteLongStr(w, f.Mechanisms); err != nil {
		return errors.New("could not write Mechanism in ConnectionStart: " + err.Error())
	}

	return nil
}

// ** ConnectionStartOk **

// Identifier returns the class ID and method ID
func (f *ConnectionStartOk) Identifier() (uint16, uint16) {
	return 10, 11
}

// MethodName returns a the name of the Method
func (f *ConnectionStartOk) MethodName() string {
	return "ConnectionStartOk"
}

// FrameType returns the frame type of the method
func (f *ConnectionStartOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ConnectionStartOk) Wait() bool {
	return true
}

func (f *ConnectionStartOk) Read(r io.Reader) (err error) {
	f.Mechanism, err = ReadShortStr(r)
	if err != nil {
		return errors.New("could not read mechanism - on ConnStartOk: " + err.Error())
	}

	f.Response, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read response string - on ConnStartOk: " + err.Error())
	}
	return nil
}

func (f *ConnectionStartOk) Write(w io.Writer) (err error) {

	if err = WriteShortStr(w, f.Mechanism); err != nil {
		return errors.New("could not write Mechanism in ConnectionStartOk: " + err.Error())
	}

	if err = WriteLongStr(w, f.Response); err != nil {
		return errors.New("could not write Response in ConnectionStartOk: " + err.Error())
	}
	return nil
}

// ** ConnectionOpen **

// Identifier returns the class ID and method ID
func (f *ConnectionOpen) Identifier() (uint16, uint16) {
	return 10, 20
}

// MethodName returns a the name of the Method
func (f *ConnectionOpen) MethodName() string {
	return "ConnectionOpen"
}

// FrameType returns the frame type of the method
func (f *ConnectionOpen) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ConnectionOpen) Wait() bool {
	return true
}

func (f *ConnectionOpen) Read(r io.Reader) (err error) {
	f.Host, err = ReadShortStr(r)
	if err != nil {
		return errors.New("could to read Host: " + err.Error())
	}
	return nil
}

func (f *ConnectionOpen) Write(w io.Writer) (err error) {

	if err = WriteShortStr(w, f.Host); err != nil {
		return errors.New("could not write Host in ConnectionOpen: " + err.Error())
	}
	return nil
}

// ** ConnectionOpenOk **

// Identifier returns the class ID and method ID
func (f *ConnectionOpenOk) Identifier() (uint16, uint16) {
	return 10, 21
}

// MethodName returns a the name of the Method
func (f *ConnectionOpenOk) MethodName() string {
	return "ConnectionOpenOk"
}

// FrameType returns the frame type of the method
func (f *ConnectionOpenOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ConnectionOpenOk) Wait() bool {
	return true
}

func (f *ConnectionOpenOk) Read(r io.Reader) (err error) {
	f.Response, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could to read response in ConnectionOpenOk: " + err.Error())
	}
	return nil
}

func (f *ConnectionOpenOk) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Response); err != nil {
		return errors.New("could not write response in ConnectionOpenOk: " + err.Error())
	}
	return
}

// ** ConnectionClose **

// Identifier returns the class ID and method ID
func (f *ConnectionClose) Identifier() (uint16, uint16) {
	return 10, 30
}

// MethodName returns a the name of the Method
func (f *ConnectionClose) MethodName() string {
	return "ConnectionClose"
}

// FrameType returns the frame type of the method
func (f *ConnectionClose) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ConnectionClose) Wait() bool {
	return true
}

func (f *ConnectionClose) Read(r io.Reader) (err error) {
	f.ReplyCode, err = ReadShort(r)
	if err != nil {
		return errors.New("could not read reply code in ConnectionClose: " + err.Error())
	}

	f.ReplyText, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read reply text in ConnectionClose: " + err.Error())
	}

	f.ClassId, err = ReadShort(r)
	if err != nil {
		return errors.New("could not read classId in ConnectionClose: " + err.Error())
	}

	f.MethodId, err = ReadShort(r)
	if err != nil {
		return errors.New("could not read MethodId in ConnectionClose: " + err.Error())
	}
	return
}

func (f *ConnectionClose) Write(w io.Writer) (err error) {

	if err = WriteShort(w, f.ReplyCode); err != nil {
		return errors.New("could not write ReplyCode in ConnectionClose: " + err.Error())
	}

	if err = WriteLongStr(w, f.ReplyText); err != nil {
		return errors.New("could not write ReplyCode in ConnectionClose: " + err.Error())
	}

	if err = WriteShort(w, f.ClassId); err != nil {
		return errors.New("could not write ClassId in ConnectionClose: " + err.Error())
	}

	if err = WriteShort(w, f.MethodId); err != nil {
		return errors.New("could not write MethodId in ConnectionClose: " + err.Error())
	}
	return nil
}

// ** ConnectionCloseOk **

// Identifier returns the class ID and method ID
func (f *ConnectionCloseOk) Identifier() (uint16, uint16) {
	return 10, 31
}

// MethodName returns a the name of the Method
func (f *ConnectionCloseOk) MethodName() string {
	return "ConnectionCloseOk"
}

// FrameType returns the frame type of the method
func (f *ConnectionCloseOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ConnectionCloseOk) Wait() bool {
	return true
}

func (f *ConnectionCloseOk) Read(r io.Reader) (err error) {
	return
}

func (f *ConnectionCloseOk) Write(w io.Writer) (err error) {

	return
}

// *******************
//    CHANNEL SPECS
// *******************

// **ChannelOpen**

// Identifier returns the class ID and method ID
func (f *ChannelOpen) Identifier() (uint16, uint16) {
	return 20, 10
}

// MethodName returns a the name of the Method
func (f *ChannelOpen) MethodName() string {
	return "ChannelOpen"
}

// FrameType returns the frame type of the method
func (f *ChannelOpen) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ChannelOpen) Wait() bool {
	return true
}

func (f *ChannelOpen) Read(r io.Reader) (err error) {
	f.Reserved, err = ReadShortStr(r)
	if err != nil {
		return errors.New("could not read reserved string in ChannelOpen: " + err.Error())
	}
	return
}

func (f *ChannelOpen) Write(w io.Writer) (err error) {

	if err = WriteShortStr(w, f.Reserved); err != nil {
		return errors.New("could not write MethodId in ConnectionClose: " + err.Error())
	}
	return nil
}

// **ChannelOpenOk**

// Identifier returns the class ID and method ID
func (f *ChannelOpenOk) Identifier() (uint16, uint16) {
	return 20, 11
}

// MethodName returns a the name of the Method
func (f *ChannelOpenOk) MethodName() string {
	return "ChannelOpenOk"
}

// FrameType returns the frame type of the method
func (f *ChannelOpenOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ChannelOpenOk) Wait() bool {
	return true
}

func (f *ChannelOpenOk) Read(r io.Reader) (err error) {
	f.Response, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read response of channel open ok: " + err.Error())
	}
	return
}

func (f *ChannelOpenOk) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Response); err != nil {
		return errors.New("could not write MethodId in ConnectionClose: " + err.Error())
	}
	return
}

// **ChannelFlow**

// Identifier returns the class ID and method ID
func (f *ChannelFlow) Identifier() (uint16, uint16) {
	return 20, 20
}

// MethodName returns a the name of the Method
func (f *ChannelFlow) MethodName() string {
	return "ChannelFlow"
}

// FrameType returns the frame type of the method
func (f *ChannelFlow) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ChannelFlow) Wait() bool {
	return true
}

func (f *ChannelFlow) Read(r io.Reader) (err error) {
	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in ChannelFlow: " + err.Error())
	}
	f.Active = (bits&(1<<0) > 0)
	return
}

func (f *ChannelFlow) Write(w io.Writer) (err error) {

	var bits byte
	if f.Active {
		bits |= 1 << 0
	}
	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write Active bit in ChannelFlow: " + err.Error())
	}
	return nil
}

// **ChannelFlowOk**

// Identifier returns the class ID and method ID
func (f *ChannelFlowOk) Identifier() (uint16, uint16) {
	return 20, 21
}

// MethodName returns a the name of the Method
func (f *ChannelFlowOk) MethodName() string {
	return "ChannelFlowOk"
}

// FrameType returns the frame type of the method
func (f *ChannelFlowOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ChannelFlowOk) Wait() bool {
	return false
}

func (f *ChannelFlowOk) Read(r io.Reader) (err error) {
	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in ChannelFlowOk: " + err.Error())
	}
	f.Active = (bits&(1<<0) > 0)
	return
}

func (f *ChannelFlowOk) Write(w io.Writer) (err error) {

	var bits byte
	if f.Active {
		bits |= 1 << 0
	}
	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write Active bit in ChannelFlowOk: " + err.Error())
	}
	return nil
}

// **ChannelClose**

// Identifier returns the class ID and method ID
func (f *ChannelClose) Identifier() (uint16, uint16) {
	return 20, 30
}

// MethodName returns a the name of the Method
func (f *ChannelClose) MethodName() string {
	return "ChannelClose"
}

// FrameType returns the frame type of the method
func (f *ChannelClose) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ChannelClose) Wait() bool {
	return true
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

func (f *ChannelClose) Write(w io.Writer) (err error) {

	if err = WriteShort(w, f.ReplyCode); err != nil {
		return errors.New("could not write ReplyCode in ChannelClose: " + err.Error())
	}

	if err = WriteLongStr(w, f.ReplyText); err != nil {
		return errors.New("could not write ReplyCode in ChannelClose: " + err.Error())
	}

	if err = WriteShort(w, f.ClassId); err != nil {
		return errors.New("could not write ClassId in ChannelClose: " + err.Error())
	}

	if err = WriteShort(w, f.MethodId); err != nil {
		return errors.New("could not write MethodId in ChannelClose: " + err.Error())
	}
	return nil
}

// **ChannelCloseOk**

// Identifier returns the class ID and method ID
func (f *ChannelCloseOk) Identifier() (uint16, uint16) {
	return 20, 31
}

// MethodName returns a the name of the Method
func (f *ChannelCloseOk) MethodName() string {
	return "ChannelCloseOk"
}

// FrameType returns the frame type of the method
func (f *ChannelCloseOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ChannelCloseOk) Wait() bool {
	return true
}

func (f *ChannelCloseOk) Read(r io.Reader) (err error) {
	return
}

func (f *ChannelCloseOk) Write(w io.Writer) (err error) {

	return
}

// *******************
//    EXCHANGE SPECS
// *******************

// **ExchangeDeclare**

// Identifier returns the class ID and method ID
func (f *ExchangeDeclare) Identifier() (uint16, uint16) {
	return 30, 10
}

// MethodName returns a the name of the Method
func (f *ExchangeDeclare) MethodName() string {
	return "ExchangeDeclare"
}

// FrameType returns the frame type of the method
func (f *ExchangeDeclare) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ExchangeDeclare) Wait() bool {
	return true && !f.NoWait
}

func (f *ExchangeDeclare) Read(r io.Reader) (err error) {
	f.Exchange, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read exchange name in ExchangeDeclare: " + err.Error())
	}

	f.Type, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read exchange type in ExchangeDeclare: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in ExchangeDeclare: " + err.Error())
	}
	f.NoWait = (bits&(1<<0) > 0)

	return
}

func (f *ExchangeDeclare) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Exchange); err != nil {
		return errors.New("could not write Exchange in ExchangeDeclare: " + err.Error())
	}

	if err = WriteLongStr(w, f.Type); err != nil {
		return errors.New("could not write Type in ExchangeDeclare: " + err.Error())
	}

	var bits byte

	if f.NoWait {
		bits |= 1 << 0
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in ExchangeDeclare: " + err.Error())
	}
	return nil
}

// **ExchangeDeclareOk**

// Identifier returns the class ID and method ID
func (f *ExchangeDeclareOk) Identifier() (uint16, uint16) {
	return 30, 11
}

// MethodName returns a the name of the Method
func (f *ExchangeDeclareOk) MethodName() string {
	return "ExchangeDeclareOk"
}

// FrameType returns the frame type of the method
func (f *ExchangeDeclareOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ExchangeDeclareOk) Wait() bool {
	return true
}

func (f *ExchangeDeclareOk) Read(r io.Reader) (err error) {
	return
}

func (f *ExchangeDeclareOk) Write(w io.Writer) (err error) {

	return
}

// **ExchangeDelete**

// Identifier returns the class ID and method ID
func (f *ExchangeDelete) Identifier() (uint16, uint16) {
	return 30, 20
}

// MethodName returns a the name of the Method
func (f *ExchangeDelete) MethodName() string {
	return "ExchangeDelete"
}

// FrameType returns the frame type of the method
func (f *ExchangeDelete) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ExchangeDelete) Wait() bool {
	return true && !f.NoWait
}

func (f *ExchangeDelete) Read(r io.Reader) (err error) {
	f.Exchange, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read exchange name in ExchangeDelete: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in ExchangeDelete: " + err.Error())
	}

	f.IfUnused = (bits&(1<<0) > 0)
	f.NoWait = (bits&(1<<1) > 0)

	return
}

func (f *ExchangeDelete) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Exchange); err != nil {
		return errors.New("could not write Exchange in ExchangeDelete: " + err.Error())
	}

	var bits byte

	if f.IfUnused {
		bits |= 1 << 0
	}

	if f.NoWait {
		bits |= 1 << 1
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in ExchangeDelete: " + err.Error())
	}
	return
}

// **ExchangeDeleteOk**

// Identifier returns the class ID and method ID
func (f *ExchangeDeleteOk) Identifier() (uint16, uint16) {
	return 30, 21
}

// MethodName returns a the name of the Method
func (f *ExchangeDeleteOk) MethodName() string {
	return "ExchangeDeleteOk"
}

// FrameType returns the frame type of the method
func (f *ExchangeDeleteOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ExchangeDeleteOk) Wait() bool {
	return true
}

func (f *ExchangeDeleteOk) Read(r io.Reader) (err error) {
	return
}

func (f *ExchangeDeleteOk) Write(w io.Writer) (err error) {

	return
}

// **ExchangeBind**

// Identifier returns the class ID and method ID
func (f *ExchangeBind) Identifier() (uint16, uint16) {
	return 30, 30
}

// MethodName returns a the name of the Method
func (f *ExchangeBind) MethodName() string {
	return "ExchangeBind"
}

// FrameType returns the frame type of the method
func (f *ExchangeBind) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ExchangeBind) Wait() bool {
	return true && !f.NoWait
}

func (f *ExchangeBind) Read(r io.Reader) (err error) {
	f.Destination, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read Destination in ExchangeBind: " + err.Error())
	}

	f.Source, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read Source in ExchangeBind: " + err.Error())
	}

	f.RoutingKey, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read RoutingKey in ExchangeBind: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in ExchangeBind: " + err.Error())
	}

	f.NoWait = (bits&(1<<0) > 0)

	return
}

func (f *ExchangeBind) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Destination); err != nil {
		return errors.New("could not write Destination in ExchangeBind: " + err.Error())
	}

	if err = WriteLongStr(w, f.Source); err != nil {
		return errors.New("could not write Source in ExchangeBind: " + err.Error())
	}

	if err = WriteLongStr(w, f.RoutingKey); err != nil {
		return errors.New("could not write Exchange in ExchangeBind: " + err.Error())
	}

	var bits byte

	if f.NoWait {
		bits |= 1 << 0
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in ExchangeBind: " + err.Error())
	}
	return
}

// **ExchangeBindOk**

// Identifier returns the class ID and method ID
func (f *ExchangeBindOk) Identifier() (uint16, uint16) {
	return 30, 31
}

// MethodName returns a the name of the Method
func (f *ExchangeBindOk) MethodName() string {
	return "ExchangeBindOk"
}

// FrameType returns the frame type of the method
func (f *ExchangeBindOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ExchangeBindOk) Wait() bool {
	return true
}

func (f *ExchangeBindOk) Read(r io.Reader) (err error) {
	return
}

func (f *ExchangeBindOk) Write(w io.Writer) (err error) {

	return
}

// **ExchangeUnbind**

// Identifier returns the class ID and method ID
func (f *ExchangeUnbind) Identifier() (uint16, uint16) {
	return 30, 40
}

// MethodName returns a the name of the Method
func (f *ExchangeUnbind) MethodName() string {
	return "ExchangeUnbind"
}

// FrameType returns the frame type of the method
func (f *ExchangeUnbind) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ExchangeUnbind) Wait() bool {
	return true && !f.NoWait
}

func (f *ExchangeUnbind) Read(r io.Reader) (err error) {
	f.Destination, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read Destination in ExchangeUnbind: " + err.Error())
	}

	f.Source, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read Source in ExchangeUnbind: " + err.Error())
	}

	f.RoutingKey, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read RoutingKey in ExchangeUnbind: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in ExchangeUnbind: " + err.Error())
	}

	f.NoWait = (bits&(1<<0) > 0)

	return
}

func (f *ExchangeUnbind) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Destination); err != nil {
		return errors.New("could not write Destination in ExchangeUnbind: " + err.Error())
	}

	if err = WriteLongStr(w, f.Source); err != nil {
		return errors.New("could not write Source in ExchangeUnbind: " + err.Error())
	}

	if err = WriteLongStr(w, f.RoutingKey); err != nil {
		return errors.New("could not write Exchange in ExchangeUnbind: " + err.Error())
	}

	var bits byte

	if f.NoWait {
		bits |= 1 << 0
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in ExchangeUnbind: " + err.Error())
	}
	return
}

// **ExchangeUnbindOk**

// Identifier returns the class ID and method ID
func (f *ExchangeUnbindOk) Identifier() (uint16, uint16) {
	return 30, 41
}

// MethodName returns a the name of the Method
func (f *ExchangeUnbindOk) MethodName() string {
	return "ExchangeUnbindOk"
}

// FrameType returns the frame type of the method
func (f *ExchangeUnbindOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *ExchangeUnbindOk) Wait() bool {
	return true
}

func (f *ExchangeUnbindOk) Read(r io.Reader) (err error) {
	return
}

func (f *ExchangeUnbindOk) Write(w io.Writer) (err error) {

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

// **QueueDeclare**

// Identifier returns the class ID and method ID
func (f *QueueDeclare) Identifier() (uint16, uint16) {
	return 40, 10
}

// MethodName returns a the name of the Method
func (f *QueueDeclare) MethodName() string {
	return "QueueDeclare"
}

// FrameType returns the frame type of the method
func (f *QueueDeclare) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *QueueDeclare) Wait() bool {
	return true && !f.NoWait
}

func (f *QueueDeclare) Read(r io.Reader) (err error) {
	f.Queue, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read queue name in QueueDeclare: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in QueueDeclare: " + err.Error())
	}
	f.NoWait = (bits&(1<<0) > 0)

	return
}

func (f *QueueDeclare) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Queue); err != nil {
		return errors.New("could not write Exchange in QueueDeclare: " + err.Error())
	}

	var bits byte

	if f.NoWait {
		bits |= 1 << 0
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in QueueDeclare: " + err.Error())
	}
	return
}

// **QueueDeclareOk**

// Identifier returns the class ID and method ID
func (f *QueueDeclareOk) Identifier() (uint16, uint16) {
	return 40, 11
}

// MethodName returns a the name of the Method
func (f *QueueDeclareOk) MethodName() string {
	return "QueueDeclareOk"
}

// FrameType returns the frame type of the method
func (f *QueueDeclareOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *QueueDeclareOk) Wait() bool {
	return true
}

func (f *QueueDeclareOk) Read(r io.Reader) (err error) {
	f.Queue, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read queue name in QueueDeclareOk: " + err.Error())
	}

	f.MessageCnt, err = ReadLong(r)
	if err != nil {
		return errors.New("could not read message count in QueueDeclareOk: " + err.Error())
	}

	f.ConsumerCnt, err = ReadLong(r)
	if err != nil {
		return errors.New("could not read message count in QueueDeclareOk: " + err.Error())
	}

	return
}

func (f *QueueDeclareOk) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Queue); err != nil {
		return errors.New("could not write Exchange in QueueDeclareOk: " + err.Error())
	}

	if err = WriteLong(w, f.MessageCnt); err != nil {
		return errors.New("could not write Exchange in QueueDeclareOk: " + err.Error())
	}

	if err = WriteLong(w, f.ConsumerCnt); err != nil {
		return errors.New("could not write Exchange in QueueDeclareOk: " + err.Error())
	}
	return

}

// **QueueBind**

// Identifier returns the class ID and method ID
func (f *QueueBind) Identifier() (uint16, uint16) {
	return 40, 20
}

// MethodName returns a the name of the Method
func (f *QueueBind) MethodName() string {
	return "QueueBind"
}

// FrameType returns the frame type of the method
func (f *QueueBind) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *QueueBind) Wait() bool {
	return true && !f.NoWait
}

func (f *QueueBind) Read(r io.Reader) (err error) {
	f.Queue, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read queue name in QueueBind: " + err.Error())
	}

	f.Exchange, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read exchange in QueueBind: " + err.Error())
	}

	f.RoutingKey, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read routingkey in QueueBind: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in QueueBind: " + err.Error())
	}
	f.NoWait = (bits&(1<<0) > 0)

	return
}

func (f *QueueBind) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Queue); err != nil {
		return errors.New("could not write Destination in QueueBind: " + err.Error())
	}

	if err = WriteLongStr(w, f.Exchange); err != nil {
		return errors.New("could not write Source in QueueBind: " + err.Error())
	}

	if err = WriteLongStr(w, f.RoutingKey); err != nil {
		return errors.New("could not write Exchange in QueueBind: " + err.Error())
	}

	var bits byte

	if f.NoWait {
		bits |= 1 << 0
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in QueueBind: " + err.Error())
	}
	return
}

// **QueueBindOk**

// Identifier returns the class ID and method ID
func (f *QueueBindOk) Identifier() (uint16, uint16) {
	return 40, 21
}

// MethodName returns a the name of the Method
func (f *QueueBindOk) MethodName() string {
	return "QueueBindOk"
}

// FrameType returns the frame type of the method
func (f *QueueBindOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *QueueBindOk) Wait() bool {
	return true
}

func (f *QueueBindOk) Read(r io.Reader) (err error) {
	return
}

func (f *QueueBindOk) Write(w io.Writer) (err error) {

	return
}

// **QueueUnbind**

// Identifier returns the class ID and method ID
func (f *QueueUnbind) Identifier() (uint16, uint16) {
	return 40, 30
}

// MethodName returns a the name of the Method
func (f *QueueUnbind) MethodName() string {
	return "QueueUnbind"
}

// FrameType returns the frame type of the method
func (f *QueueUnbind) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *QueueUnbind) Wait() bool {
	return true
}

func (f *QueueUnbind) Read(r io.Reader) (err error) {
	f.Queue, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read queue name in QueueUnbind: " + err.Error())
	}

	f.Exchange, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read exchange in QueueUnbind: " + err.Error())
	}

	f.RoutingKey, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read routingkey in QueueUnbind: " + err.Error())
	}

	return
}

func (f *QueueUnbind) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Queue); err != nil {
		return errors.New("could not write Destination in QueueUnbind: " + err.Error())
	}

	if err = WriteLongStr(w, f.Exchange); err != nil {
		return errors.New("could not write Source in QueueUnbind: " + err.Error())
	}

	if err = WriteLongStr(w, f.RoutingKey); err != nil {
		return errors.New("could not write Exchange in QueueUnbind: " + err.Error())
	}
	return
}

// **QueueUnbindOk**

// Identifier returns the class ID and method ID
func (f *QueueUnbindOk) Identifier() (uint16, uint16) {
	return 40, 31
}

// MethodName returns a the name of the Method
func (f *QueueUnbindOk) MethodName() string {
	return "QueueUnbindOk"
}

// FrameType returns the frame type of the method
func (f *QueueUnbindOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *QueueUnbindOk) Wait() bool {
	return true
}

func (f *QueueUnbindOk) Read(r io.Reader) (err error) {
	return
}

func (f *QueueUnbindOk) Write(w io.Writer) (err error) {

	return
}

// **QueueDelete**

// Identifier returns the class ID and method ID
func (f *QueueDelete) Identifier() (uint16, uint16) {
	return 40, 40
}

// MethodName returns a the name of the Method
func (f *QueueDelete) MethodName() string {
	return "QueueDelete"
}

// FrameType returns the frame type of the method
func (f *QueueDelete) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *QueueDelete) Wait() bool {
	return true && !f.NoWait
}

func (f *QueueDelete) Read(r io.Reader) (err error) {
	f.Queue, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read queue name in QueueDelete: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in QueueDelete: " + err.Error())
	}
	f.IfUnused = (bits&(1<<0) > 0)
	f.IfEmpty = (bits&(1<<1) > 0)
	f.NoWait = (bits&(1<<2) > 0)

	return
}

func (f *QueueDelete) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Queue); err != nil {
		return errors.New("could not write Queue in QueueDelete: " + err.Error())
	}

	var bits byte

	if f.IfUnused {
		bits |= 1 << 0
	}
	if f.IfEmpty {
		bits |= 1 << 1
	}
	if f.NoWait {
		bits |= 1 << 2
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in QueueDelete: " + err.Error())
	}
	return
}

// **QueueDeleteOk**

// Identifier returns the class ID and method ID
func (f *QueueDeleteOk) Identifier() (uint16, uint16) {
	return 40, 41
}

// MethodName returns a the name of the Method
func (f *QueueDeleteOk) MethodName() string {
	return "QueueDeleteOk"
}

// FrameType returns the frame type of the method
func (f *QueueDeleteOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *QueueDeleteOk) Wait() bool {
	return true
}

func (f *QueueDeleteOk) Read(r io.Reader) (err error) {
	f.MessageCnt, err = ReadLong(r)
	if err != nil {
		return errors.New("could not read queue name in QueueDeleteOk: " + err.Error())
	}

	return
}

func (f *QueueDeleteOk) Write(w io.Writer) (err error) {

	if err = WriteLong(w, f.MessageCnt); err != nil {
		return errors.New("could not write MessageCnt in QueueDeleteOk: " + err.Error())
	}
	return
}

// *******************
//    QoS SPECS
//        basicConsume - 10
//        basicConsumeOk - 11
//        basicCancel - 20
//        basicCancelOk - 21
//        basicPublish - 30
//        basicReturn  - 40
//        basicDeliver - 50
//        basicAck     - 60
//        basicNack    - 70
// *******************

// ** BasicConsume **

// Identifier returns the class ID and method ID
func (f *BasicConsume) Identifier() (uint16, uint16) {
	return 50, 10
}

// MethodName returns a the name of the Method
func (f *BasicConsume) MethodName() string {
	return "BasicConsume"
}

// FrameType returns the frame type of the method
func (f *BasicConsume) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *BasicConsume) Wait() bool {
	return true && !f.NoWait
}

func (f *BasicConsume) Read(r io.Reader) (err error) {
	f.Queue, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read queue name in basicConsume: " + err.Error())
	}

	f.ConsumerTag, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read consumer tag in basicConsume: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in basicConsume: " + err.Error())
	}

	f.NoAck = (bits&(1<<0) > 0)
	f.NoWait = (bits&(1<<1) > 0)

	return
}

func (f *BasicConsume) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Queue); err != nil {
		return errors.New("could not write Queue in BasicConsume: " + err.Error())
	}

	if err = WriteLongStr(w, f.ConsumerTag); err != nil {
		return errors.New("could not write ConsumerTag in BasicConsume: " + err.Error())
	}

	var bits byte

	if f.NoAck {
		bits |= 1 << 0
	}
	if f.NoWait {
		bits |= 1 << 1
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in BasicConsume: " + err.Error())
	}
	return
}

// **BasicConsumeOk**

// Identifier returns the class ID and method ID
func (f *BasicConsumeOk) Identifier() (uint16, uint16) {
	return 50, 11
}

// MethodName returns a the name of the Method
func (f *BasicConsumeOk) MethodName() string {
	return "BasicConsumeOk"
}

// FrameType returns the frame type of the method
func (f *BasicConsumeOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *BasicConsumeOk) Wait() bool {
	return true
}

func (f *BasicConsumeOk) Read(r io.Reader) (err error) {
	f.ConsumerTag, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read consumer tag in basicConsumeOk: " + err.Error())
	}

	return
}

func (f *BasicConsumeOk) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.ConsumerTag); err != nil {
		return errors.New("could not write ConsumerTag in BasicConsumeOk: " + err.Error())
	}
	return
}

// **BasicCancel**'

// Identifier returns the class ID and method ID
func (f *BasicCancel) Identifier() (uint16, uint16) {
	return 50, 20
}

// MethodName returns a the name of the Method
func (f *BasicCancel) MethodName() string {
	return "BasicCancel"
}

// FrameType returns the frame type of the method
func (f *BasicCancel) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *BasicCancel) Wait() bool {
	return true
}

func (f *BasicCancel) Read(r io.Reader) (err error) {
	f.ConsumerTag, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read consumer tag in BasicCancel: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in BasicCancel: " + err.Error())
	}

	f.NoWait = (bits&(1<<0) > 0)

	return
}

func (f *BasicCancel) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.ConsumerTag); err != nil {
		return errors.New("could not write ConsumerTag in BasicCancel: " + err.Error())
	}

	var bits byte
	if f.NoWait {
		bits |= 1 << 0
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in BasicCancel: " + err.Error())
	}
	return
}

// **BasicCancelOk**

// Identifier returns the class ID and method ID
func (f *BasicCancelOk) Identifier() (uint16, uint16) {
	return 50, 21
}

// MethodName returns a the name of the Method
func (f *BasicCancelOk) MethodName() string {
	return "BasicCancelOk"
}

// FrameType returns the frame type of the method
func (f *BasicCancelOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *BasicCancelOk) Wait() bool {
	return true
}

func (f *BasicCancelOk) Read(r io.Reader) (err error) {
	f.ConsumerTag, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read consumer tag in BasicCancelOk: " + err.Error())
	}

	return
}

func (f *BasicCancelOk) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.ConsumerTag); err != nil {
		return errors.New("could not write ConsumerTag in BasicCancelOk: " + err.Error())
	}
	return
}

// ** BasicPublish **

// Identifier returns the class ID and method ID
func (f *BasicPublish) Identifier() (uint16, uint16) {
	return 50, 30
}

// MethodName returns a the name of the Method
func (f *BasicPublish) MethodName() string {
	return "BasicPublish"
}

// FrameType returns the frame type of the method
func (f *BasicPublish) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *BasicPublish) Wait() bool {
	return false
}

// GetBody gets the method frame body
func (f *BasicPublish) GetBody() []byte {
	return f.Body
}

// SetBody sets the method frame body
func (f *BasicPublish) SetBody(b []byte) {
	f.Body = b
}

func (f *BasicPublish) Read(r io.Reader) (err error) {
	f.Exchange, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read exchange name in BasicPublish: " + err.Error())
	}

	f.RoutingKey, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read consumer tag in BasicPublish: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in BasicPublish: " + err.Error())
	}

	f.Immediate = (bits&(1<<0) > 0)

	return
}

func (f *BasicPublish) Write(w io.Writer) (err error) {

	if err = WriteLongStr(w, f.Exchange); err != nil {
		return errors.New("could not write Exchange in BasicPublish: " + err.Error())
	}

	if err = WriteLongStr(w, f.RoutingKey); err != nil {
		return errors.New("could not write RoutingKey in BasicPublish: " + err.Error())
	}

	var bits byte

	if f.Immediate {
		bits |= 1 << 0
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in BasicPublish: " + err.Error())
	}
	return
}

// ** BasicReturn **

// Identifier returns the class ID and method ID
func (f *BasicReturn) Identifier() (uint16, uint16) {
	return 50, 40
}

// MethodName returns a the name of the Method
func (f *BasicReturn) MethodName() string {
	return "BasicReturn"
}

// FrameType returns the frame type of the method
func (f *BasicReturn) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *BasicReturn) Wait() bool {
	return false
}

// GetBody gets the method frame body
func (f *BasicReturn) GetBody() []byte {
	return f.Body
}

// SetBody sets the method frame body
func (f *BasicReturn) SetBody(b []byte) {
	f.Body = b
}

func (f *BasicReturn) Read(r io.Reader) (err error) {
	f.ReplyCode, err = ReadShort(r)
	if err != nil {
		return errors.New("could not read reply code in BasicReturn: " + err.Error())
	}

	f.ReplyText, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read reply test in BasicReturn: " + err.Error())
	}

	f.Exchange, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read exchange name in BasicReturn: " + err.Error())
	}

	f.RoutingKey, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read consumer tag in BasicReturn: " + err.Error())
	}

	return
}

func (f *BasicReturn) Write(w io.Writer) (err error) {

	if err = WriteShort(w, f.ReplyCode); err != nil {
		return errors.New("could not write ReplyCode in BasicReturn: " + err.Error())
	}

	if err = WriteLongStr(w, f.ReplyText); err != nil {
		return errors.New("could not write ReplyText in BasicReturn: " + err.Error())
	}

	if err = WriteLongStr(w, f.Exchange); err != nil {
		return errors.New("could not write Exchange in BasicReturn: " + err.Error())
	}

	if err = WriteLongStr(w, f.RoutingKey); err != nil {
		return errors.New("could not write RoutingKey in BasicReturn: " + err.Error())
	}
	return nil
}

// ** BasicDeliver **

// Identifier returns the class ID and method ID
func (f *BasicDeliver) Identifier() (uint16, uint16) {
	return 50, 50
}

// MethodName returns a the name of the Method
func (f *BasicDeliver) MethodName() string {
	return "BasicDeliver"
}

// FrameType returns the frame type of the method
func (f *BasicDeliver) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *BasicDeliver) Wait() bool {
	return false
}

// GetBody gets the method frame body
func (f *BasicDeliver) GetBody() []byte {
	return f.Body
}

// SetBody sets the method frame body
func (f *BasicDeliver) SetBody(b []byte) {
	f.Body = b
}

func (f *BasicDeliver) Read(r io.Reader) (err error) {
	f.ConsumerTag, err = ReadShortStr(r)
	if err != nil {
		return errors.New("could not read consumer tag name in BasicDeliver: " + err.Error())
	}

	f.DeliveryTag, err = ReadLongLong(r)
	if err != nil {
		return errors.New("could not read exchange name in BasicDeliver: " + err.Error())
	}

	f.Exchange, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read exchange name in BasicDeliver: " + err.Error())
	}

	f.RoutingKey, err = ReadLongStr(r)
	if err != nil {
		return errors.New("could not read consumer tag in BasicDeliver: " + err.Error())
	}

	return
}

func (f *BasicDeliver) Write(w io.Writer) (err error) {

	if err = WriteShortStr(w, f.ConsumerTag); err != nil {
		return errors.New("could not write ConsumerTag in BasicDeliver: " + err.Error())
	}

	if err = WriteLongLong(w, f.DeliveryTag); err != nil {
		return errors.New("could not write ReplyText in BasicDeliver: " + err.Error())
	}

	if err = WriteLongStr(w, f.Exchange); err != nil {
		return errors.New("could not write Exchange in BasicDeliver: " + err.Error())
	}

	if err = WriteLongStr(w, f.RoutingKey); err != nil {
		return errors.New("could not write RoutingKey in BasicDeliver: " + err.Error())
	}
	return
}

// ** BasicAck **

// Identifier returns the class ID and method ID
func (f *BasicAck) Identifier() (uint16, uint16) {
	return 50, 60
}

// MethodName returns a the name of the Method
func (f *BasicAck) MethodName() string {
	return "BasicAck"
}

// FrameType returns the frame type of the method
func (f *BasicAck) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *BasicAck) Wait() bool {
	return false
}

func (f *BasicAck) Read(r io.Reader) (err error) {
	f.DeliveryTag, err = ReadLongLong(r)
	if err != nil {
		return errors.New("could not read exchange name in BasicAck: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in BasicAck: " + err.Error())
	}

	f.Multiple = (bits&(1<<0) > 0)

	return
}

func (f *BasicAck) Write(w io.Writer) (err error) {

	if err = WriteLongLong(w, f.DeliveryTag); err != nil {
		return errors.New("could not write DeliveryTag in BasicAck: " + err.Error())
	}

	var bits byte
	if f.Multiple {
		bits |= 1 << 0
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in BasicAck: " + err.Error())
	}
	return
}

// ** BasicNack **

// Identifier returns the class ID and method ID
func (f *BasicNack) Identifier() (uint16, uint16) {
	return 50, 70
}

// MethodName returns a the name of the Method
func (f *BasicNack) MethodName() string {
	return "BasicNack"
}

// FrameType returns the frame type of the method
func (f *BasicNack) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *BasicNack) Wait() bool {
	return false
}

func (f *BasicNack) Read(r io.Reader) (err error) {

	f.DeliveryTag, err = ReadLongLong(r)
	if err != nil {
		return errors.New("could not read exchange name in BasicNack: " + err.Error())
	}

	bits, err := ReadOctet(r)
	if err != nil {
		return errors.New("could not read bits in BasicNack: " + err.Error())
	}

	f.Multiple = (bits&(1<<0) > 0)
	f.Requeue = (bits&(1<<1) > 0)
	return
}

func (f *BasicNack) Write(w io.Writer) (err error) {

	if err = WriteLongLong(w, f.DeliveryTag); err != nil {
		return errors.New("could not write DeliveryTag in BasicNack: " + err.Error())
	}

	var bits byte
	if f.Multiple {
		bits |= 1 << 0
	}
	if f.Requeue {
		bits |= 1 << 1
	}

	if err = WriteOctet(w, bits); err != nil {
		return errors.New("could not write bits in BasicNack: " + err.Error())
	}
	return
}

// *******************
//   Tx SPECS
//   Class - 60
//	 TxSelect - 10
//	 TxSelectOk - 11
//	 TxCommit - 20
//	 TxCommitOk - 21
//	 TxRollback - 30
//	 TxRollbackOk - 31
// *******************

// ** TxSelect **

// Identifier returns the class ID and method ID
func (f *TxSelect) Identifier() (uint16, uint16) {
	return 60, 10
}

// MethodName returns a the name of the Method
func (f *TxSelect) MethodName() string {
	return "TxSelect"
}

// FrameType returns the frame type of the method
func (f *TxSelect) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *TxSelect) Wait() bool {
	return true
}

func (f *TxSelect) Read(r io.Reader) (err error) {
	return
}

func (f *TxSelect) Write(w io.Writer) (err error) {
	return
}

// ** TxSelectOk **

// Identifier returns the class ID and method ID
func (f *TxSelectOk) Identifier() (uint16, uint16) {
	return 60, 11
}

// MethodName returns a the name of the Method
func (f *TxSelectOk) MethodName() string {
	return "TxSelectOk"
}

// FrameType returns the frame type of the method
func (f *TxSelectOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *TxSelectOk) Wait() bool {
	return true
}

func (f *TxSelectOk) Read(r io.Reader) (err error) {
	return
}

func (f *TxSelectOk) Write(w io.Writer) (err error) {
	return
}

// ** TxCommit **

// Identifier returns the class ID and method ID
func (f *TxCommit) Identifier() (uint16, uint16) {
	return 60, 20
}

// MethodName returns a the name of the Method
func (f *TxCommit) MethodName() string {
	return "TxCommit"
}

// FrameType returns the frame type of the method
func (f *TxCommit) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *TxCommit) Wait() bool {
	return true
}

func (f *TxCommit) Read(r io.Reader) (err error) {
	return
}

func (f *TxCommit) Write(w io.Writer) (err error) {
	return
}

// ** TxCommitOk **

// Identifier returns the class ID and method ID
func (f *TxCommitOk) Identifier() (uint16, uint16) {
	return 60, 21
}

// MethodName returns a the name of the Method
func (f *TxCommitOk) MethodName() string {
	return "TxCommitOk"
}

// FrameType returns the frame type of the method
func (f *TxCommitOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *TxCommitOk) Wait() bool {
	return true
}

func (f *TxCommitOk) Read(r io.Reader) (err error) {
	return
}

func (f *TxCommitOk) Write(w io.Writer) (err error) {
	return
}

// ** TxRollback **

// Identifier returns the class ID and method ID
func (f *TxRollback) Identifier() (uint16, uint16) {
	return 60, 30
}

// MethodName returns a the name of the Method
func (f *TxRollback) MethodName() string {
	return "TxRollback"
}

// FrameType returns the frame type of the method
func (f *TxRollback) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *TxRollback) Wait() bool {
	return true
}

func (f *TxRollback) Read(r io.Reader) (err error) {
	return
}

func (f *TxRollback) Write(w io.Writer) (err error) {
	return
}

// ** TxRollbackOk **

// Identifier returns the class ID and method ID
func (f *TxRollbackOk) Identifier() (uint16, uint16) {
	return 60, 31
}

// MethodName returns a the name of the Method
func (f *TxRollbackOk) MethodName() string {
	return "TxRollbackOk"
}

// FrameType returns the frame type of the method
func (f *TxRollbackOk) FrameType() byte {
	return 1
}

// Wait returns a boolean signifying if the method need any wait
func (f *TxRollbackOk) Wait() bool {
	return true
}

func (f *TxRollbackOk) Read(r io.Reader) (err error) {
	return
}

func (f *TxRollbackOk) Write(w io.Writer) (err error) {
	return
}
