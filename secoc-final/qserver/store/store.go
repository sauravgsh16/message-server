package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sauravgsh16/secoc-third/secoc-final/proto"
)

var CONTENT_BUCKET = []byte("content")
var INDEX_BUCKET = []byte("index")

type Key struct {
	id        int64
	queuename string
}

type MsgStore struct {
	db          *bolt.DB
	index       map[int64]*proto.IndexMessage
	messages    map[int64]*proto.Message
	qmToAdd     map[Key]*proto.QueueMessage
	qmToDelete  map[Key]*proto.QueueMessage
	qmDelivered map[Key]*proto.QueueMessage
	indexMux    sync.RWMutex
	msgMux      sync.RWMutex
	persistMux  sync.Mutex
}

func deleteFileIfPresent(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}
	if err := os.Remove(filePath); err != nil {
		panic("Failed to remove old DB file")
	}
}

func New(filePath string) (*MsgStore, error) {
	// Check if file name already present -
	// remove it for new session to start
	deleteFileIfPresent(filePath)

	db, err := bolt.Open(filePath, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &MsgStore{
		db:          db,
		index:       make(map[int64]*proto.IndexMessage),
		messages:    make(map[int64]*proto.Message),
		qmToAdd:     make(map[Key]*proto.QueueMessage),
		qmToDelete:  make(map[Key]*proto.QueueMessage),
		qmDelivered: make(map[Key]*proto.QueueMessage),
	}, nil
}

func (ms *MsgStore) Start() {
	go ms.handlePeriodicPersists()
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

func (ms *MsgStore) GetMessage(id int64) (*proto.Message, bool) {
	ms.msgMux.Lock()
	defer ms.msgMux.Unlock()

	msg, found := ms.messages[id]
	return msg, found
}

func (ms *MsgStore) GetIndexMessage(id int64) (*proto.IndexMessage, bool) {
	ms.indexMux.Lock()
	defer ms.indexMux.Unlock()

	im, found := ms.index[id]
	return im, found
}

func (ms *MsgStore) RemoveRef(qm *proto.QueueMessage, queueName string, mrh []proto.MessageResourceHolder) error {

	im, found := ms.GetIndexMessage(qm.ID)
	if !found {
		panic("Message in queue - but not in Index. Unrecoverable failure")
	}

	if len(queueName) == 0 {
		panic("Bad Queue Name")
	}

	im.Refs -= 1
	if im.Refs == 0 {
		ms.msgMux.Lock()
		delete(ms.messages, qm.ID)
		ms.msgMux.Unlock()

		ms.indexMux.Lock()
		delete(ms.index, qm.ID)
		ms.indexMux.Unlock()
	}

	for _, rh := range mrh {
		rh.ReleaseResources(qm)
	}
	return nil
}

func (ms *MsgStore) Get(qm *proto.QueueMessage, mrh []proto.MessageResourceHolder) (*proto.Message, bool) {
	ms.msgMux.Lock()
	defer ms.msgMux.Unlock()

	resources := make([]proto.MessageResourceHolder, 0, len(mrh))
	for _, rh := range mrh {
		if !rh.AcquireResources(qm) {
			break
		}
		resources = append(resources, rh)
	}

	if len(resources) == len(mrh) {
		msg, found := ms.messages[qm.ID]
		if !found {
			panic("Message not Found.")
		}
		return msg, true
	}

	// Failure, release resource
	for _, rh := range resources {
		rh.ReleaseResources(qm)
	}
	return nil, false
}

func calcMessageSize(msg *proto.Message) uint32 {
	return uint32(len(msg.Payload))
}

func (ms *MsgStore) handlePeriodicPersists() {
	interval := time.Duration(200 * time.Millisecond)
	for {
		time.Sleep(interval)
		ms.persistDB()
	}
}

func (ms *MsgStore) resetOps() {
	ms.qmToAdd = make(map[Key]*proto.QueueMessage)
	ms.qmToDelete = make(map[Key]*proto.QueueMessage)
	ms.qmDelivered = make(map[Key]*proto.QueueMessage)
}

func (ms *MsgStore) persistDB() {
	ms.persistMux.Lock()
	qmToAdd := ms.qmToAdd
	qmToDelete := ms.qmToDelete
	qmDelivered := ms.qmDelivered
	ms.resetOps()
	ms.persistMux.Unlock()

	toDelete := make([]Key, 0, len(qmToAdd))

	for id, _ := range qmToDelete {
		if _, ok := qmToAdd[id]; ok {
			delete(qmToAdd, id)
			toDelete = append(toDelete, id)
		}
		delete(qmDelivered, id)
	}
	for _, id := range toDelete {
		delete(qmToDelete, id)
	}

	// Update db to persist new changes
	uf := ms.updateFunc(qmToAdd, qmToDelete, qmDelivered)
	if err := ms.db.Update(uf); err != nil {
		panic("Failed to persist data: " + err.Error())
	}
}

func (ms *MsgStore) updateFunc(qmToAdd, qmToDelete, qmDelivered map[Key]*proto.QueueMessage) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		// Add functionality
		alreadyAdded := make(map[int64]bool)
		for k, qm := range qmToAdd {
			// Add messages to content/index stores
			if _, ok := alreadyAdded[k.id]; !ok {
				msg, foundMsg := ms.GetMessage(k.id)
				im, foundIm := ms.GetIndexMessage(k.id)
				if foundMsg != foundIm {
					panic("Message index discrrpency")
				}
				if !foundMsg {
					// This means, msg must have been deleted earlier
					continue
				}
				persistMessage(tx, msg)
				persistIndexMessage(tx, im)
			}
			persistQueueMessage(tx, k.queuename, qm)
		}

		// Update delivered
		for k, qm := range qmDelivered {
			persistQueueMessage(tx, k.queuename, qm)
		}

		// Delete qm - remove from queue
		for k, qm := range qmToDelete {
			if err := depersistQueueMessage(tx, k.queuename, qm.ID); err != nil {
				return err
			}
			refCount, err := decrementIndexRef(tx, qm.ID, ms)
			if err != nil {
				return err
			}

			// Delete messages if no references remain
			if refCount == 0 {
				if err := depersistMessage(tx, qm.ID); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func getIdByte(id int64) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(buf, binary.LittleEndian, id)
	return buf.Bytes()
}

func persistMessage(tx *bolt.Tx, msg *proto.Message) error {
	bucket, err := tx.CreateBucketIfNotExists(CONTENT_BUCKET)
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	key := getIdByte(msg.ID)
	return bucket.Put(key, encoded)
}

func persistIndexMessage(tx *bolt.Tx, im *proto.IndexMessage) error {
	bucket, err := tx.CreateBucketIfNotExists(CONTENT_BUCKET)
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(im)
	if err != nil {
		return err
	}
	key := getIdByte(im.ID)
	return bucket.Put(key, encoded)
}

func persistQueueMessage(tx *bolt.Tx, qname string, qm *proto.QueueMessage) error {
	bucketName := fmt.Sprintf("queue_%s", qname)
	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(qm)
	if err != nil {
		return err
	}
	key := getIdByte(qm.ID)
	return bucket.Put(key, encoded)
}

func depersistMessage(tx *bolt.Tx, id int64) error {
	bucket, err := tx.CreateBucketIfNotExists(CONTENT_BUCKET)
	if err != nil {
		return err
	}
	key := getIdByte(id)
	return bucket.Delete(key)
}

func depersistQueueMessage(tx *bolt.Tx, qname string, id int64) error {
	bucketName := fmt.Sprintf("queue_%s", qname)
	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return err
	}
	key := getIdByte(id)
	return bucket.Delete(key)
}

func decrementIndexRef(tx *bolt.Tx, id int64, ms *MsgStore) (int32, error) {
	bucket, err := tx.CreateBucketIfNotExists(INDEX_BUCKET)
	if err != nil {
		return -1, err
	}

	im := &proto.IndexMessage{}

	key := getIdByte(id)
	dataBytes := bucket.Get(key)

	err = json.Unmarshal(dataBytes, im)
	if err != nil {
		return -1, err
	}

	if im.Refs < 1 {
		panic("Index messages reference count - negative")
	}

	im.Refs -= 1
	if im.Refs == 0 {
		// Reference count Zero - delete ref from messages and Index.
		// Remove key from bucket
		ms.msgMux.Lock()
		delete(ms.messages, id)
		ms.msgMux.Unlock()

		ms.indexMux.Lock()
		delete(ms.index, id)
		ms.indexMux.Unlock()
		return 0, bucket.Delete(key)
	}

	freshEncodedBytes, err := json.Marshal(im)
	if err != nil {
		return -1, nil
	}
	return im.Refs, bucket.Put(key, freshEncodedBytes)
}
