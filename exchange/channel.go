package exchange

import (
        "sync"
)

type Channel struct {
        connection *Connection // see if this is required
        ID         int
        m          sync.Mutex
        rpc        chan message
        Exchanges  map[string]*exchangeDeclare
}

func newChannel(c *Connection, id int) *Channel {
        ch := &Channel{
                connection: c,
                rpc:        make(chan message), // See if this is required
                ID:         id,
                Exchanges:  make(map[string]*exchangeDeclare),
        }
        return ch
}

func (ch *Channel) DeclareExchange(name, extype string) {
        ex := &exchangeDeclare{
                Exchange: name,
                Type:     extype,  // Implement validator of type
        }
        ch.Exchanges[name] = ex
}

/*
Firstly, whenever we connect to Rabbit we need a fresh, empty queue.
To do this we could create a queue with a random name,
or, even better - let the server choose a random queue name for us.

Secondly, once we disconnect the consumer the queue should be automatically deleted.
*/