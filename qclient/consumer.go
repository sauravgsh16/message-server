package qclient

import "sync"

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

func (c *Consumers) send(consumer string, d *Delivery) {
	c.mux.Lock()
	defer c.mux.Unlock()

	ch, found := c.consumerMap[consumer]
	if found {
		ch <- d
	}
}

func (c *Consumers) add(consumertag string, dChan chan Delivery) {
	c.mux.Lock()
	defer c.mux.Unlock()

	// if found, close the previous channel for conflict resolution
	if _, ok := c.consumerMap[consumertag]; ok {
		close(dChan)
	}

	in := make(chan *Delivery)
	c.consumerMap[consumertag] = in

	c.wg.Add(1)
	go c.createBufferOrConsume(in, dChan)
}

func (c *Consumers) createBufferOrConsume(in chan *Delivery, dChan chan Delivery) {
	defer c.wg.Done()
	defer close(dChan)

	inLooper := in
	buf := make([]*Delivery, 0)

	for delivery := range in {
		buf = append(buf, delivery)

		for len(buf) > 0 {
			select {
			case <-c.closed:
				return
			case delivery, stillconsuming := <-inLooper:
				if stillconsuming {
					buf = append(buf, delivery)
				} else {
					inLooper = nil
				}
			case dChan <- *buf[0]:
				buf = buf[1:]
			}
		}
	}
}

func (c *Consumers) cancel(consumer string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	ch, found := c.consumerMap[consumer]
	if found {
		delete(c.consumerMap, consumer)
		close(ch)
	}

	return found
}

func (c *Consumers) close() {
	c.mux.Lock()
	defer c.mux.Unlock()

	close(c.closed)

	for tag, ch := range c.consumerMap {
		delete(c.consumerMap, tag)
		close(ch)
	}

	// Wait till we get a done from all the goroutines called.
	c.wg.Wait()
}
