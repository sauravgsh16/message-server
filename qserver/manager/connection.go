package manager

import (
        "fmt"
        "io"
        "net"
        "sync"

        "github.com/sauravgsh16/secoc-third/constant"
        "github.com/sauravgsh16/secoc-third/allocate"
)

type ConnectionManager struct {
        wg         sync.WaitGroup
        mux        sync.Mutex
        Listener   net.Listener
        destructor sync.Once                  // to close all queue, if connection is closed
        clients    map[int]*client // To change to io.ReadWriteClose, implement close method on client
                                              // Can have a close channel, which will indicate closing of all in queues handled
                                              // by the client
        allocator  *allocate.Allocator
}

func NewConnection() *ConnectionManager {
        c := &ConnectionManager{
                clients:   make(map[int]*client),
                allocator: allocate.NewAllocator(),
        }
        addr, err := net.ResolveTCPAddr("tcp", constant.UnsecuredPort)
        if err != nil {
                fmt.Println("invalid port") // Add logging
                panic("unable to resolve TCP address")
        }

        l, err := net.ListenTCP("tcp", addr)
        if err != nil {
                fmt.Printf("Error listening on TCP socket: %s\n", err.Error())
                panic(err.Error())
        }
        fmt.Printf("TCP listening on port %s\n", constant.UnsecuredPort)
        c.Listener = l
        return c
}

func (cm *ConnectionManager) Start() {
        cm.wg.Add(1)
        go cm.acceptTCPConn()
        cm.wg.Wait()
}

func (cm *ConnectionManager) acceptTCPConn() {
        defer cm.wg.Done()
        for {
                conn, err := cm.Listener.Accept()
                if err != nil {
                        // NEEDS ERROR HANDLING
                        fmt.Printf("error occurred whist opening TCP connection: %s\n", err.Error())
                        conn.Close()
                        continue
                }
                cm.mux.Lock()
                id, ok := cm.allocator.Next()
                if !ok {
                        // NEEDS ERROR HANDLING
                        fmt.Println("max number of connections")
                        return
                }
                nc := newClient(id, conn)
                cm.clients[id] = nc
                cm.mux.Unlock()
                cm.wg.Add(1)
                fmt.Printf("Accepted connection, client: %d, conn: %v\n", id, conn)
                go cm.handleConnection(id)
        }
}

func (cm *ConnectionManager) handleConnection(id int) {
        buf := make([]byte, buffersize)
        cl, ok := cm.clients[id]
        if !ok {
                fmt.Printf("no client with id: %d present\n", id)
                return
        }
        for {
                size, err := cl.r.ReadByte()
                if err != nil {
                        fmt.Printf("Error reading data: %s", err.Error())
                }
                _, err = io.ReadFull(cl.r, buf[:int(size)])
                if err != nil {
                        fmt.Printf("Error reading all the bytes from connection, %s", err.Error())
                }
                fmt.Printf("len:%d  data:%v", len(buf), buf)
                break
        }
}
