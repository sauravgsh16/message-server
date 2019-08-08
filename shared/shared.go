package shared


// Message will be the struct created after reading the bytes on the
// connection. This will be sent to the queue. The queue will store it
// accordingly.
type Message struct {
	Body []byte
}

func (m Message) IsEmpty() bool {
	return len(m.Body) == 0
}