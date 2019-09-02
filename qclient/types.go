package qclient

type Confirmation struct {
	DeliveryTag uint64
	State       bool
}
