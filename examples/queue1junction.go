package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
	"time"
)

type Queue1[T any] struct {
	head *QueueElement1[T]
}

type QueueElement1[T any] struct {
	valueSignal   func(types.Unit) (T, error)
	setNextSignal func(QueueElement1[T])
	getNextSignal func(types.Unit) (QueueElement1[T], error)
}

func newQueue1[T any]() (func(T), func(types.Unit) (T, error)) {
	j := junction.NewJunction()

	firstPort, firstSignal := junction.NewAsyncSignal[QueueElement1[T]](j)
	lastPort, lastSignal := junction.NewAsyncSignal[QueueElement1[T]](j)

	enqueuePort, enqueueSignal := junction.NewAsyncSignal[T](j)
	dequeuePort, dequeueSignal := junction.NewSyncSignal[types.Unit, T](j)

	junction.NewBinaryAsyncJoinPattern[QueueElement1[T], T](lastPort, enqueuePort).Action(func(last QueueElement1[T], value T) {
		elem := newQueueElement1[T](j, value)

		lastSignal(elem)

		if last.valueSignal != nil {
			last.setNextSignal(elem)
		} else {
			firstSignal(elem)
		}
	})

	junction.NewBinarySyncJoinPattern[QueueElement1[T], types.Unit, T](firstPort, dequeuePort).Action(func(first QueueElement1[T], a types.Unit) T {
		nextSignal, _ := first.getNextSignal(types.Unit{})

		firstSignal(nextSignal)

		val, _ := first.valueSignal(types.Unit{})

		return val
	})

	lastSignal(QueueElement1[T]{})

	return enqueueSignal, dequeueSignal
}

func newQueueElement1[T any](j *junction.Junction, value T) QueueElement1[T] {
	valuePort, valueSignal := junction.NewSyncSignal[types.Unit, T](j)
	setNextPort, setNextSignal := junction.NewAsyncSignal[QueueElement1[T]](j)
	getNextPort, getNextSignal := junction.NewSyncSignal[types.Unit, QueueElement1[T]](j)

	junction.NewUnarySyncJoinPattern[types.Unit, T](valuePort).Action(func(a types.Unit) T {
		return value
	})

	junction.NewBinarySyncJoinPattern[types.Unit, QueueElement1[T], QueueElement1[T]](getNextPort, setNextPort).
		Action(func(a types.Unit, node QueueElement1[T]) QueueElement1[T] {
			return node
		})

	return QueueElement1[T]{
		valueSignal:   valueSignal,
		setNextSignal: setNextSignal,
		getNextSignal: getNextSignal,
	}
}

func main() {
	enqueue, dequeue := newQueue1[int]()

	// Producing items
	go func() {
		for i := 0; ; i++ {
			go func(num int) {
				fmt.Println("Enqueueing ", num)
				enqueue(num)
			}(i)
			time.Sleep(time.Second)
		}
	}()

	// Consuming items
	go func() {
		for i := 0; ; i++ {
			go func() {
				val, _ := dequeue(types.Unit{})
				fmt.Println("Dequeued: ", val)
			}()
			time.Sleep(time.Second)
		}
	}()

	for true {
	}
}
