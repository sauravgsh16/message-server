package main

import (
        "fmt"
        "net/rpc"
        "github.com/sauravgsh16/secoc-third/shared"
)

func main() {
        c, err := rpc.Dial("tcp", ":1234")
        if err != nil {
                fmt.Println(err)
                return
        }
        defer c.Close()

        a1 := &shared.Value{A:2, B:3}
        var r1 int

        err = c.Call("V.Add", a1, &r1)
        if err != nil {
                fmt.Println(err)
                return
        }
        fmt.Printf("Reply r1: %d\n", r1)

        a2 := &shared.Name{N:"foo"}
        err = c.Call("N.AddNum", a2, &r1)
        if err != nil {
                fmt.Println("HERE", err)
                return
        }

        r2 := []string{}
        err = c.Call("N.GetName", a2, &r2)
        if err != nil {
                fmt.Printf("NOW HERE %v\n", err)
        }

        fmt.Printf("Reply r2: %+v\n", r2)
}