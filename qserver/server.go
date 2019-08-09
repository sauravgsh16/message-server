package main

import (
        "github.com/sauravgsh16/secoc-third/qserver/manager"
)

func main() {
        c := manager.NewConnection()
        c.Start()
}

