package shared


// QData : need to check if more implementation details required
// Check Table : - types.go amqp
type QData struct {
	data string
}

type QCreate struct {
	Qtype	   string
	QName      string
	AutoDelete bool   // Delete Queue automatically after all info sent
}



// NOT USED
type Value struct {
	A, B int
}

type Valuer interface {
	Add(args *Value, reply *int) error
}

type Name struct {
	N string
}

type Namer interface {
	AddNum(args Name, reply *int) error
	GetName(_ interface{}, reply []string) error
}