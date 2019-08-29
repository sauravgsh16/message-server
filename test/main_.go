package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type Message struct {
	Name string
	Age  int
}

func openDB(filePath string) *bolt.DB {
	db, err := bolt.Open(filePath, 0666, nil)
	if err != nil {
		log.Fatalf("DB open failed")
	}
	return db
}

func persist(db *bolt.DB, msg *Message) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("msg"))
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		return b.Put([]byte(msg.Name), encoded)
	})
	return err
}

func get(db *bolt.DB) (*Message, error) {
	msg := &Message{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("msg")).Get([]byte("Test"))
		if err := json.Unmarshal(b, msg); err != nil {
			return err
		}
		return nil
	})
	return msg, err
}

func main_() {
	msg := &Message{Name: "Test", Age: 16}
	db := openDB("store.db")
	if err := persist(db, msg); err != nil {
		fmt.Printf(err.Error())
	}
	m, err := get(db)
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf("%v\n", m)
}
