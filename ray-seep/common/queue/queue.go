package queue

import (
	"errors"
	"sync"
)

type elem struct {
	value    interface{}
	previous *elem
	Next     *elem
}

type Queue interface {
	Pop() *elem
	Push(*elem) error
	Shift() *elem
	Length() int64
}

type queue struct {
	lock     sync.Mutex
	head     *elem
	rear     *elem
	length   int64
	capacity int64
}

func NewQueue(args ...interface{}) *queue {
	max := int64(100000000)
	if len(args) > 0 {
		switch args[0].(type) {
		case int64:
			max = args[0].(int64)
		case int:
			max = int64(args[0].(int))
		}
	}
	return &queue{
		capacity: max,
	}
}
func (sel *queue) Pop() *elem {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if sel.length <= 0 {
		sel.length = 0
		return nil
	}
	val := sel.head
	sel.head = sel.head.Next
	sel.length--
	val.previous = nil
	val.Next = nil
	return val
}
func (sel *queue) Shift() *elem {
	sel.lock.Lock()
	defer sel.lock.Unlock()

	if sel.length <= 0 {
		sel.length = 0
		return nil
	}

	val := sel.rear
	if sel.rear.previous == nil {
		sel.rear = sel.head
	} else {
		sel.rear = sel.rear.previous
		sel.rear.Next = nil
	}
	val.previous = nil
	val.Next = nil
	sel.length--
	return val
}
func (sel *queue) Push(n *elem) error {
	if sel.length >= sel.capacity {
		return errors.New("over max num for stack")
	}
	sel.push(n)
	return nil
}
func (sel *queue) push(top *elem) {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if 0 == sel.length {
		sel.head = top
		sel.rear = sel.head
	}
	top.Next = sel.head
	sel.head.previous = top
	sel.head = top
	sel.length++
	return
}
func (sel *queue) Length() int64 {
	return sel.length
}
