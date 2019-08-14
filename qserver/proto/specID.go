package proto

import "errors"

type Spec struct {
        ClassID  uint16
        Methods  map[string]uint16
}

func NewSpec(classID uint16) Spec {
        return Spec{ClassID: classID}
}

type SpecMap map[string]Spec

func (sm SpecMap) addClass(id uint16, class string) {
        sm[class] = NewSpec(10)
}

func (sm SpecMap) addMethod(methodId uint16, clsName, methodName string) error {
        spec, ok := sm[clsName]
        if !ok {
                return errors.New("Spec not present: " + clsName)
        }
        spec.Methods[methodName] = methodId
        return nil
}

var s SpecMap

func init() {
        classes := []string{"connection", "channel"}
        classInit(classes)
        methodsInit()
}

func classInit(classes []string) {
        classId := uint16(10)
        for _, cls := range classes {
                s.addClass(classId, cls)
                classId += 10
        }
}

func methodsInit() {

}

/*
connection - 10 class
    method:
        connectionStart    - 10
        connectionStartOk  - 11
        connectionSecure   - 20
        connectionSecureOk - 21
        connectionTune     - 30
        connectionTuneOk   - 31
        connectionOpen     - 40
        connectionOpenOk   - 41
        connectionClose    - 50
        connectionCloseOk  - 51
        connectionBlocked  - 60
        connectionUnblocked- 61
channel - 20 class
    method:
        channelOpen - 10
        channelOpenOk - 11
        channelFlow - 20
        channelFlowOk - 21
        channelClose - 40
        channelCloseOk - 41
*/


