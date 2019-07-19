package queue

import (
	"errors"
        "fmt"
        
        sh "github.com/sauravgsh16/secoc-third/shared"
)

type msg struct {
	next  *msg
	value sh.QData
}

// List - linked
type List struct {
	Root *msg
	len  int
	ch   chan sh.QData
}

// Len of list
func (l *List) Len() int {
	return l.len
}

// NewList points to pointer to a new list
func newList() *List {
	l := &List{
		ch: make(chan sh.QData),
	}
	return l
}

func (l *List) findLast() *msg {
	cur := l.Root
	for cur.next != nil {
		cur = cur.next
	}
	return cur
}

func (l *List) append(d sh.QData) {
	n := &msg{value: d}
	if l.Root == nil {
		l.Root = n
		return
	}
	last := l.findLast()
	last.next = n
}

func (l *List) remove() (sh.QData, error) {
	if l.Root == nil {
		return sh.QData{}, errors.New("Cannot remove from empty list")
	}
	n := *l.Root
	l.Root = n.next
	n.next = nil // remove reference to next pointer
	l.len--
	return n.value, nil
}

// Append to end of list
func (l *List) Append(d sh.QData) {
	l.append(d)
}

// Remove one msg
func (l *List) Remove() sh.QData {
	d, err := l.remove()
	if err != nil {
		return sh.QData{}
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