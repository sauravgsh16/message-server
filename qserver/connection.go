package qserver

import (
        "fmt"
        "io"
        "net"
        "sync"

        "github.com/sauravgsh16/secoc-third/constant"
        "github.com/sauravgsh16/secoc-third/allocate"
)

type Connection struct {
        wg         sync.WaitGroup
        mux        sync.Mutex
        Listener   net.Listener
        destructor sync.Once             // to close all queue, if connection is closed
        clients    map[int]io.ReadWriter // To change to io.ReadWriteClose, implement close method on client
                                         // Can have a close channel, which will indicate closing of all in queues handled
                                         // by the client
        allocator  *allocate.Allocator
}

func NewConnection() *Connection {
        c := &Connection{
                clients:   make(map[int]io.ReadWriter),
                allocator: allocate.NewAllocator(),
        }
        addr, err := net.ResolveTCPAddr("tcp", constant.UnsecuredPort)
        if err != nil {
                fmt.Println("invalid port") // Add logging
                panic("unable to resolve TCP address")
        }

        l, err := net.ListenTCP("tcp", addr)
        if err != nil {
                fmt.Printf("Error listening on TCP socket: %s", err.Error())
                panic(err.Error())
        }
        c.Listener = l
        return c
}

func (c *Connection) Start() {
        c.wg.Add(1)
        go c.acceptTCPConn()
        c.wg.Wait()
}

func (c *Connection) acceptTCPConn() {
        defer c.wg.Done()
        for {
                conn, err := c.Listener.Accept()
                if err != nil {
                        // NEEDS ERROR HANDLING
                        fmt.Printf("error occurred whist opening TCP connection: %s", err.Error())
                }
                c.mux.Lock()
                id, ok := c.allocator.Next()
                if !ok {
                        // NEEDS ERROR HANDLING
                        fmt.Printf("max number of connections")
                        return
                }
                nc := NewClient(id, conn)
                c.clients[id] = nc
                c.mux.Unlock()
                c.wg.Add(1)
                go c.handleClient(id)
        }
}

func (c *Connection) handleClient(id int) {
        
}
