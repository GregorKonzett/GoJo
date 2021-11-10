package main

import (
	"../gojo/junction"
	"../gojo/types"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Queue1[T any] struct {
	head *QueueElement1[T]
	j    *junction.Junction
}

type QueueElement1[T any] struct {
	get  func(types.Unit) (T, error)
	next *QueueElement1[T]
}

func (q *Queue1[T]) enqueue(val T) {

	getPort, get := junction.NewSyncSignal[types.Unit, T]((*q).j)
	setPort, set := junction.NewAsyncSignal[T]((*q).j)

	junction.NewBinarySyncJoinPattern[types.Unit, T, T](getPort, setPort).Action(func(a types.Unit, val T) T {
		return val
	})

	set(val)

	queueElem := QueueElement1[T]{
		get: get,
	}

	if (*q).head == nil {
		(*q).head = &queueElem
	} else {
		cur := (*q).head
		for cur.next != nil {
			cur = cur.next
		}

		cur.next = &queueElem
	}
}

func (q *Queue1[T]) dequeue() (T, error) {
	if (*q).head != nil {
		val, _ := (*q).head.get(types.Unit{})
		(*q).head = (*q).head.next
		return val, nil
	} else {
		var returnData T
		return returnData, errors.New("empty queue")
	}
}

func main() {
	queue := Queue1[int]{j: junction.NewJunction()}
	mutex := sync.Mutex{}
	size := 10

	for i := 0; i < size; i++ {
		fmt.Println("Enqueueing ", i)
		queue.enqueue(i)
	}

	for i := 0; i < size; i++ {
		go func() {
			time.Sleep(time.Second)
			// Figure out how to do without locks
			mutex.Lock()
			val, _ := queue.dequeue()
			mutex.Unlock()
			fmt.Println("Dequeued: ", val)
		}()
	}

	for true {
	}
}
