package main

import (
        "fmt"
        "log"
        "net"
        "net/rpc"
        "sync"

        "github.com/sauravgsh16/secoc-third/server/queue"
        sh "github.com/sauravgsh16/secoc-third/shared"
)

const (
        PORT = "9001"
        INIT_QUEUE_ID = 1
)

type queueref  map[string]*queue.Queue

func (qr queueref) add(n string, q *queue.Queue) {
        qr[n] = q
}

func (qr queueref) get(n string) (*queue.Queue, error) {
        q, ok := qr[n]
        if ! ok {
                return &queue.Queue{}, fmt.Errorf("queue with id: %s - not found", n)
        }
        return q, nil
}

type queuetype map[string]queueref

type ServerStore struct {
        qNum    int     // Initialized with INIT_QUEUE_ID, keep incrementing for new queue
        queues  queuetype
        wmux    sync.Mutex
        rwmux   sync.RWMutex
}

// registers new queue, returns id of queue which has been created
func (ss *ServerStore) newQ(qtype, qName string) {
        q := queue.NewQueue(ss.qNum)
        qr := ss.queues[qtype]
        ss.wmux.Lock()
        qr.add(qName, q)
        ss.wmux.Unlock()
}

func (ss *ServerStore) getQueue(qtype string , id int) (*queue.Queue, error) {
        switch qtype {
        case "broadcast":
                ss.rwmux.Lock()
                q, ok := ss.bqueues[id]
                ss.rwmux.Unlock()
                if !ok {
                        return &queue.Queue{}, fmt.Errorf("queue with id: %d - not found", id)
                }
                return q, nil
        case "send":
                ss.rwmux.Lock()
                q, ok := ss.squeues[id]
                ss.rwmux.Unlock()
                if !ok {
                        return &queue.Queue{}, fmt.Errorf("queue with id: %d - not found", id)
                }
                return q, nil
        default:
                return &queue.Queue{}, fmt.Errorf("invalid queue type %v", qtype)
        }
}

type HandleQ struct {
        ss *ServerStore
}

func (h *HandleQ) CreateQueue(args *sh.QCreate, res *int) error {
        // res will contain id of queue which was created
        switch args.Qtype {
        case "broadcast", "send":
                h.ss.newQ(args.Qtype, args.QName)
                *res = 1
        default:
                return fmt.Errorf("invalid Queue type: %s", args.Qtype)
        }
        return nil
}



func registerhandlers() error {
        return nil
}

func main() {
        if err := registerhandlers(); err != nil {
                log.Fatalf("Failed to register handlers %v", err)
        }
        l, err := net.Listen("tcp", fmt.Sprintf(":%s", PORT))
        if err != nil {
                log.Fatalf("Unable to start Queue server %v", err)
        }
        for {
                conn, err := l.Accept()
                if err != nil {
                        log.Fatalf("Unable to accept connection on listener %v", err)
                }
                go rpc.ServeConn(conn)
        }
}