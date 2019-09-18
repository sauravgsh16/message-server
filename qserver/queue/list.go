package queue

import (
	"errors"
	"fmt"
	"sync"
)

type qData interface{}

type msg struct {
	next  *msg
	value qData
}

// List - linked
type List struct {
	Root *msg
	len  int
	mux  sync.Mutex
}

// Len of list
func (l *List) Len() int {
	return l.len
}

// NewList points to pointer to a new list
func newlist() *List { return &List{} }

func (l *List) removeRef() *List {
	l.Root = &msg{}
	l.len = 0
	return l
}

func (l *List) findLast() *msg {
	cur := l.Root
	for cur.next != nil {
		cur = cur.next
	}
	return cur
}

func (l *List) append(d qData) {
	n := &msg{value: d}
	if l.Root == nil {
		l.Root = n
		return
	}
	last := l.findLast()
	last.next = n
}

func (l *List) remove() error {
	if l.Root == nil {
		return errors.New("Cannot remove from empty list")
	}
	n := *l.Root
	l.Root = n.next
	n.next = nil // remove reference to next pointer
	l.len--
	return nil
}

// Append to end of list
func (l *List) Append(d qData) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.append(d)
	l.len++
}

// Remove one msg
func (l *List) Remove() {
	l.mux.Lock()
	defer l.mux.Unlock()

	if err := l.remove(); err != nil {
		panic(err.Error())
	}
}

func (l *List) Front() qData {
	if l.len == 0 {
		return nil
	}
	return l.Root.value
}

func (l *List) String() string {
	data := ""
	cur := l.Root
	for cur != nil {
		val := fmt.Sprintf("%v->", cur.value)
		data = data + val
		cur = cur.next
	}
	return fmt.Sprintf("%s, of lenght %d", data, l.len)
}
