package exchange


type payload struct {
        msg []byte
}

type dataStore struct {
        data     payload
        receiver []string // need to change this- needs to be []queues
                          // also needs to capability to mark if data has been
                          // sent to the queue.
}