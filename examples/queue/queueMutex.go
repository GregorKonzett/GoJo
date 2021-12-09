package queue

import (
	"../../gojo/types"
	"sync"
	"time"
)

type QueueElementMutex[T any] struct {
	val  T
	next *QueueElementMutex[T]
}

type QueueMutex[T any] struct {
	head *QueueElementMutex[T]
	tail *QueueElementMutex[T]
}

func NewQueueMutex[T any]() (func(T), func(types.Unit) (T, error)) {
	m := sync.Mutex{}
	c := sync.NewCond(&m)

	queue := QueueMutex[T]{}

	enqueue := func(val T) {
		c.L.Lock()

		elem := &QueueElementMutex[T]{
			val: val,
		}

		if queue.head == nil {
			queue.head = elem
			queue.tail = elem
		} else {
			queue.tail.next = elem
			queue.tail = elem
		}

		c.Signal()
		c.L.Unlock()
	}

	dequeue := func(types.Unit) (T, error) {
		c.L.Lock()
		for queue.head == nil {
			c.Wait()
		}

		firstElem := queue.head
		queue.head = queue.head.next
		time.Sleep(time.Duration(10))
		c.L.Unlock()

		return (*firstElem).val, nil
	}

	return enqueue, dequeue
}
