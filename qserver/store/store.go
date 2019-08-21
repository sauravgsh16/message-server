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
	messages   map[int64]*proto.Message
	addOps     map[Key]*proto.QueueMessage
	delOps     map[Key]*proto.QueueMessage
	deliverOps map[Key]*proto.QueueMessage
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

func (s *MsgStore) Persist() {

}

func (s *MsgStore) AddMessage(msg *proto.Message, qs []string) (map[string][]*proto.QueueMessage, error) {
}

func (s *MsgStore) RemoveRef(qm *proto.QueueMessage, queue string) {}
