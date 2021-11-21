package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
	"math/rand"
	"time"
)

type QueueElement1[T any] struct {
	getValueSignal func(types.Unit) (T, error)
	setNextSignal  func(QueueElement1[T])
	getNextSignal  func(types.Unit) (QueueElement1[T], error)
}

func newQueue1[T any]() (func(T), func(types.Unit) (T, error)) {
	j := junction.NewJunction()

	firstPort, firstSignal := junction.NewAsyncSignal[QueueElement1[T]](j)
	lastPort, lastSignal := junction.NewAsyncSignal[QueueElement1[T]](j)
	emptyLastPort, emptyLastSignal := junction.NewSyncSignal[types.Unit, types.Unit](j)

	enqueuePort, enqueueSignal := junction.NewAsyncSignal[T](j)
	dequeuePort, dequeueSignal := junction.NewSyncSignal[types.Unit, T](j)

	tail := newQueueElement1[T](j, QueueElement1[T]{})

	junction.NewBinarySyncJoinPattern[QueueElement1[T], types.Unit, types.Unit](lastPort, emptyLastPort).Action(func(last QueueElement1[T], a types.Unit) types.Unit {
		return types.Unit{}
	})

	junction.NewBinaryAsyncJoinPattern[QueueElement1[T], T](lastPort, enqueuePort).Action(func(last QueueElement1[T], value T) {
		elem := newQueueElement1[T](j, tail)
		insertIntoElement1[T](j, &elem, value)

		lastSignal(elem)

		if last.getValueSignal != nil {
			last.getNextSignal(types.Unit{})
			last.setNextSignal(elem)
		} else {
			firstSignal(elem)
		}
	})

	junction.NewBinarySyncJoinPattern[QueueElement1[T], types.Unit, T](firstPort, dequeuePort).Action(func(first QueueElement1[T], a types.Unit) T {
		nextSignal, _ := first.getNextSignal(types.Unit{})

		if nextSignal.getValueSignal != nil {
			firstSignal(nextSignal)
		} else {
			emptyLastSignal(types.Unit{})
			lastSignal(tail)
		}

		val, _ := first.getValueSignal(types.Unit{})
		return val
	})

	lastSignal(tail)

	return enqueueSignal, dequeueSignal
}

func newQueueElement1[T any](j *junction.Junction, tail QueueElement1[T]) QueueElement1[T] {
	setNextPort, setNextSignal := junction.NewAsyncSignal[QueueElement1[T]](j)
	getNextPort, getNextSignal := junction.NewSyncSignal[types.Unit, QueueElement1[T]](j)

	junction.NewBinarySyncJoinPattern[types.Unit, QueueElement1[T], QueueElement1[T]](getNextPort, setNextPort).
		Action(func(a types.Unit, node QueueElement1[T]) QueueElement1[T] {
			return node
		})

	setNextSignal(tail)

	return QueueElement1[T]{
		setNextSignal: setNextSignal,
		getNextSignal: getNextSignal,
	}
}

func insertIntoElement1[T any](j *junction.Junction, elem *QueueElement1[T], val T) {
	getValuePort, getValueSignal := junction.NewSyncSignal[types.Unit, T](j)

	junction.NewUnarySyncJoinPattern[types.Unit, T](getValuePort).
		Action(func(a types.Unit) T {
			return val
		})

	(*elem).getValueSignal = getValueSignal
}

func main() {
	enqueue, dequeue := newQueue1[int]()
	producerCount := 2
	consumerCount := 2

	// Producer
	for i := 1; i <= producerCount; i++ {
		func(num int) {
			val := num
			for val < 5 {
				time.Sleep(time.Second * time.Duration(rand.Intn(2)))
				fmt.Println("Producer", num, " Enqueueing ", val)
				enqueue(val)
				val += 1 * num
			}
		}(i)
	}

	// Consumer
	for i := 0; i < consumerCount; i++ {
		go func(num int) {
			for true {
				time.Sleep(time.Second * time.Duration(rand.Intn(2)))
				val, _ := dequeue(types.Unit{})
				fmt.Println("Consumer", num, " consuming ", val)
			}
		}(i)
	}

	for true {
	}
}
