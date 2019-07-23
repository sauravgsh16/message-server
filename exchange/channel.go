package exchange

import (
        "sync"
)

type Channel struct {
        connection *Connection // see if this is required
        ID         int
        m          sync.Mutex
        rpc        chan message
}

func newChannel(c *Connection, id int) *Channel {
        ch := &Channel{
                connection: c,
                rpc:        make(chan message), // See if this is required
                ID:         id,
        }
        return ch
}

func (ch *Channel) DeclareExchange() {
        
}