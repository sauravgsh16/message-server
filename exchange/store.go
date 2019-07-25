package exchange


type data struct {
        msg []byte
}

type dataStore struct {
        kind     string
        exchange string   // name of exchange
        data     data
        receiver []string // need to change this- needs to be []queues
                          // also needs to capability to mark if data has been
                          // sent to the queue.
}