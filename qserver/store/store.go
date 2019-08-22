package store

import (
	"sync"

	"github.com/boltdb/bolt"
	"github.com/sauravgsh16/secoc-third/qserver/proto"
)

type Key struct {
	id        int64
	queuename string
}

type MsgStore struct {
	db         *bolt.DB
	index      map[int64]*proto.IndexMessage
	messages   map[int64]*proto.Message
	addOps     map[Key]*proto.QueueMessage
	delOps     map[Key]*proto.QueueMessage
	deliverOps map[Key]*proto.QueueMessage
	indexMux   sync.RWMutex
	msgMux     sync.RWMutex
	persistMux sync.Mutex
}

func New(db *bolt.DB) (*MsgStore, error) {
	// Check if file name already present -
	// remove it for new session to start
	return &MsgStore{
		db:         db,
		messages:   make(map[int64]*proto.Message),
		addOps:     make(map[Key]*proto.QueueMessage),
		delOps:     make(map[Key]*proto.QueueMessage),
		deliverOps: make(map[Key]*proto.QueueMessage),
	}, nil
}

func (ms *MsgStore) Persist() {

}

func (ms *MsgStore) AddMessage(msg *proto.Message, qs []string) (map[string][]*proto.QueueMessage, error) {
	msgs := make([]*proto.TxMessage, 0, len(qs))
	for _, q := range qs {
		msgs = append(msgs, proto.NewTxMessage(msg, q))
	}
	return ms.AddTxMessages(msgs)
}

func (ms *MsgStore) AddTxMessages(msgs []*proto.TxMessage) (map[string][]*proto.QueueMessage, error) {
	// Create IndexMessage instances for each message
	indexMessage := make(map[int64]*proto.IndexMessage)
	queueMessage := make(map[string][]*proto.QueueMessage)

	for _, msg := range msgs {
		// Check or Create index message
		im, found := indexMessage[msg.Msg.ID]
		if !found {
			im = proto.NewIndexMessage(msg.Msg.ID, 0, 0)
			indexMessage[msg.Msg.ID] = im
		}
		im.Refs += 1

		// Check or Create Queue Message
		queues, found := queueMessage[msg.QueueName]
		if !found {
			queues = make([]*proto.QueueMessage, 0, 1)
		}
		qm := proto.NewQueueMessage(
			msg.Msg.ID,
			0,
			calcMessageSize(msg.Msg),
		)
		queueMessage[msg.QueueName] = append(queues, qm)
	}

	// Add indexes and messages to Memory
	ms.msgMux.Lock()
	defer ms.msgMux.Unlock()
	ms.indexMux.Lock()
	defer ms.indexMux.Unlock()

	for _, msg := range msgs {
		ms.index[msg.Msg.ID] = indexMessage[msg.Msg.ID]
		ms.messages[msg.Msg.ID] = msg.Msg
	}
	return queueMessage, nil
}

func calcMessageSize(msg *proto.Message) uint32 {
	size := uint32(0)
	for _, frame := range msg.Payload {
		size += uint32(len(frame.Payload))
	}
	return size
}

func (ms *MsgStore) GetIndex(id int64) (*proto.IndexMessage, bool) {
	ms.indexMux.Lock()
	defer ms.indexMux.Unlock()

	im, found := ms.index[id]
	return im, found
}

func (ms *MsgStore) RemoveRef(qm *proto.QueueMessage, queueName string, mrh []proto.MessageResourceHolder) error {

	im, found := ms.GetIndex(qm.ID)
	if !found {
		panic("Message in queue - but not in Index. Unrecoverable failure")
	}

	if len(queueName) == 0 {
		panic("Bad Queue Name")
	}

	_ = im
	return nil
}
