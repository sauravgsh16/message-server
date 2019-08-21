package queue

import (
	"errors"
	"fmt"
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

func (l *List) remove() (qData, error) {
	if l.Root == nil {
		return nil, errors.New("Cannot remove from empty list")
	}
	n := *l.Root
	l.Root = n.next
	n.next = nil // remove reference to next pointer
	l.len--
	return n.value, nil
}

// Append to end of list
func (l *List) Append(d qData) {
	l.append(d)
}

// Remove one msg
func (l *List) Remove() qData {
	d, err := l.remove()
	if err != nil {
		// return qData{data: sh.Message{Body: make([]byte, 0)}}
		return nil
	}
	return d
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
