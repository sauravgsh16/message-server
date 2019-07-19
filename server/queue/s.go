package queue

import (
        "fmt"
        "log"
        "net/rpc"
        "net"
        _ "sync"
        "github.com/sauravgsh16/secoc-third/shared"
)

type V struct {}

type N struct {
        name []string
}

func (v *V) Add(args *shared.Value, reply *int) error {
        *reply = args.A + args.B
        return nil
}

func (n *N) AddNum(args *shared.Name, reply *int) error {
        n.name = append(n.name, args.N)
        fmt.Printf("%+v", n)
        *reply = 0
        return nil
}

func (n *N) GetName(args *shared.Name, reply *[]string,) error {
        *reply = append(*reply, n.name...)
        return nil
}

func tempmain() {
        myV := new(V)
        rpc.Register(myV)
        rpc.Register(new(N))
        /*
        t, err := net.ResolveTCPAddr("tcp4", ":1234")
        if err != nil {
                fmt.Println(err)
                return
        }
        l, err := net.ListenTCP("tcp4", t)
        if err != nil {
                fmt.Println(err)
                return
        }

        for {
                c, err := l.Accept()
                if err != nil {
                        continue
                }
                fmt.Printf("%s\n", c.RemoteAddr())
                rpc.ServeConn(c)
        }
        */
        // l, e := net.Listen("("tcp", fmt.Sprintf(":%(":%v", ":1234"))
        l, e := net.Listen("tcp", fmt.Sprintf(":%s", "1234"))
        if e != nil {
                log.Fatal("listen error:", e)
        }
        // http.Serve(l, nil)
        for {
                conn, _ := l.Accept()
                go rpc.ServeConn(conn)
        } 
}