package queue

import (
	"../../gojo/junction"
	"../../gojo/types"
	"fmt"
	"sync"
	"time"
)

type QueueElement1[T any] struct {
	getValueSignal func(types.Unit) (T, error)
	setValueSignal func(T)
	setNextSignal  func(QueueElement1[T])
	getNextSignal  func(types.Unit) (QueueElement1[T], error)
}

func NewQueue1[T any]() (func(T), func(types.Unit) (T, error)) {
	j := junction.NewJunction()

	firstPort, firstSignal := junction.NewAsyncPort[QueueElement1[T]](j)
	lastPort, lastSignal := junction.NewAsyncPort[QueueElement1[T]](j)

	enqueuePort, enqueueSignal := junction.NewAsyncPort[T](j)
	dequeuePort, dequeueSignal := junction.NewSyncPort[types.Unit, T](j)

	junction.NewBinaryAsyncJoinPattern[QueueElement1[T], T](lastPort, enqueuePort).Action(func(last QueueElement1[T], value T) {
		elem := newQueueElement1[T](j)

		last.setValueSignal(value)
		last.setNextSignal(elem)

		lastSignal(elem)
	})

	junction.NewBinarySyncJoinPattern[QueueElement1[T], types.Unit, T](firstPort, dequeuePort).Action(func(first QueueElement1[T], a types.Unit) T {
		nextSignal, _ := first.getNextSignal(types.Unit{})
		val, _ := first.getValueSignal(types.Unit{})

		firstSignal(nextSignal)

		return val
	})

	elem := newQueueElement1[T](j)

	firstSignal(elem)
	lastSignal(elem)

	return enqueueSignal, dequeueSignal
}

func newQueueElement1[T any](j *junction.Junction) QueueElement1[T] {
	setNextPort, setNextSignal := junction.NewAsyncPort[QueueElement1[T]](j)
	getNextPort, getNextSignal := junction.NewSyncPort[types.Unit, QueueElement1[T]](j)
	getValuePort, getValueSignal := junction.NewSyncPort[types.Unit, T](j)
	setValuePort, setValueSignal := junction.NewAsyncPort[T](j)

	junction.NewBinarySyncJoinPattern[types.Unit, QueueElement1[T], QueueElement1[T]](getNextPort, setNextPort).
		Action(func(a types.Unit, node QueueElement1[T]) QueueElement1[T] {
			return node
		})

	junction.NewBinarySyncJoinPattern[T, types.Unit, T](setValuePort, getValuePort).
		Action(func(val T, a types.Unit) T {
			return val
		})

	return QueueElement1[T]{
		setNextSignal:  setNextSignal,
		getNextSignal:  getNextSignal,
		setValueSignal: setValueSignal,
		getValueSignal: getValueSignal,
	}
}

func main1() {
	enqueue, dequeue := NewQueue1[int]()
	producerCount := 1000000
	consumerCount := 900

	var wg sync.WaitGroup
	start := time.Now()

	// Producer
	for i := 0; i < producerCount; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			time.Sleep(500)
			//fmt.Println("Producer", num, " Enqueueing ", num)
			enqueue(num)
		}(i)
	}

	// Consumer
	for i := 0; i < consumerCount; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()

			time.Sleep(500)
			dequeue(types.Unit{})
		}(i)
	}

	wg.Wait()
	end := time.Since(start)
	fmt.Println("Duration: ", end)
}
