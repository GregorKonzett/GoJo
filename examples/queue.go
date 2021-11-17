package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
	"time"
)

type QueueElement[T any] struct {
	valueSignal      func(types.Unit) (T, error)
	setNextSignal    func(QueueElement[T])
	getNextSignal    func(types.Unit) (QueueElement[T], error)
	hasNextSignal    func(types.Unit) (bool, error)
	setHasNextSignal func(bool)
}

func newQueue[T any]() (func(T), func(types.Unit) (T, error)) {
	j := junction.NewJunction()

	firstPort, firstSignal := junction.NewAsyncSignal[QueueElement[T]](j)
	lastPort, lastSignal := junction.NewAsyncSignal[QueueElement[T]](j)

	enqueuePort, enqueueSignal := junction.NewAsyncSignal[T](j)
	dequeuePort, dequeueSignal := junction.NewSyncSignal[types.Unit, T](j)

	junction.NewBinaryAsyncJoinPattern[QueueElement[T], T](lastPort, enqueuePort).Action(func(last QueueElement[T], value T) {
		elem := newQueueElement[T](value)

		lastSignal(elem)

		if last.valueSignal != nil {
			last.setNextSignal(elem)
		} else {
			firstSignal(elem)
		}
	})

	junction.NewBinarySyncJoinPattern[QueueElement[T], types.Unit, T](firstPort, dequeuePort).Action(func(first QueueElement[T], a types.Unit) T {
		nextSignal, _ := first.getNextSignal(types.Unit{})
		firstSignal(nextSignal)

		val, _ := first.valueSignal(types.Unit{})
		return val
	})

	//firstSignal(QueueElement[T]{})
	lastSignal(QueueElement[T]{})

	return enqueueSignal, dequeueSignal
}

func newQueueElement[T any](value T) QueueElement[T] {
	j := junction.NewJunction()

	valuePort, valueSignal := junction.NewSyncSignal[types.Unit, T](j)
	setNextPort, setNextSignal := junction.NewAsyncSignal[QueueElement[T]](j)
	getNextPort, getNextSignal := junction.NewSyncSignal[types.Unit, QueueElement[T]](j)

	junction.NewUnarySyncJoinPattern[types.Unit, T](valuePort).Action(func(a types.Unit) T {
		return value
	})

	junction.NewBinarySyncJoinPattern[types.Unit, QueueElement[T], QueueElement[T]](getNextPort, setNextPort).
		Action(func(a types.Unit, node QueueElement[T]) QueueElement[T] {
			return node
		})

	return QueueElement[T]{
		valueSignal:   valueSignal,
		setNextSignal: setNextSignal,
		getNextSignal: getNextSignal,
	}
}

func main() {
	enqueue, dequeue := newQueue[int]()

	// Producing items
	func() {
		for i := 0; i < 10; i++ {
			go func(num int) {

				fmt.Println("Enqueueing ", num)
				enqueue(num)
			}(i)
			time.Sleep(time.Second)
		}
	}()

	// Consuming items
	func() {
		for i := 0; i < 10; i++ {
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
