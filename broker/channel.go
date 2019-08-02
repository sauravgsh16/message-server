package broker

import (
        "sync"
        "sync/atomic"
)

type Channel struct {
        connection *Connection // see if this is required
        ID         uint16
        mux        sync.Mutex
        rpc        chan message
        Exchanges  map[string]*exchangeDeclare
        closed     int32
}


func newChannel(c *Connection, id uint16) *Channel {
        ch := &Channel{
                connection: c,
                rpc:        make(chan message), // See if this is required
                ID:         id,
                Exchanges:  make(map[string]*exchangeDeclare),
        }
        return ch
}

func (ch *Channel) shutdown(err *Error) {
        // TODO
}

func (ch *Channel) send(msg message) error {
        if atomic.LoadInt32(&ch.closed) == 1 {
                return ch.sendClosed(msg)
        }
        return ch.sendOpen(msg)
}

func (ch *Channel) sendClosed(msg message) error {
        if _, ok := msg.(*channelCloseOk); ok {
                return ch.connection.send(&methodFrame{
                        ChannelId: ch.ID,
                        Method:    msg,
                })
        }
        return ErrClosed
}

func (ch *Channel) sendOpen(msg message) error {
        if mwc, ok := msg.(messageWithContent); ok {
                props, body := mwc.getContent()

                err := ch.connection.send(&methodFrame{
                        ChannelId: ch.ID,
                        Method:    mwc,
                })
                if err != nil {
                        return err
                }

                err = ch.connection.send(&headerFrame{
                        ChannelId:  ch.ID,
                        Size:       uint64(len(body)),
                        Properties: props,
                })
                if err != nil {
                        return err
                }

                err = ch.connection.send(&bodyFrame{
                        ChannelId: ch.ID,
                        Body:      body,
                })
                if err != nil {
                        return err
                }
        } else {
                err := ch.connection.send(&methodFrame{
                        ChannelId: ch.ID,
                        Method:    msg,
                })
                if err != nil {
                        return err
                }
        }
        return nil
}

func (ch *Channel) Publish(exchange, key string, msg Publishing) error {
        ch.mux.Lock()
        ch.mux.Unlock()

        err := ch.send(&basicPublish{
                Exchange:   exchange,
                RoutingKey: key,
                Body:       msg.Body,
                Properties: properties{
                        UserId:    msg.UserId,
                        MessageId: msg.MessageId,
                },
        })
        if err != nil {
                return err
        }
        return nil
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