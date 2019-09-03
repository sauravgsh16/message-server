package qclient

import "sync"

// TO MOVE TO RELEVANT PLACE
// ****************************************
type Delivery struct{}

// *****************************************

type Consumers struct {
	wg          sync.WaitGroup
	closed      chan struct{}
	mux         sync.Mutex
	consumerMap map[string]chan *Delivery
}

func CreateNewConsumers() *Consumers {
	return &Consumers{
		closed:      make(chan struct{}),
		consumerMap: make(map[string]chan *Delivery),
	}
}

func (c *Consumers) add(consumer string, dChan chan Delivery) {

}

func (c *Consumers) cancel(consumer string) {

}
