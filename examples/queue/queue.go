package queue

import (
	"../../gojo/junction"
	"../../gojo/types"
	"fmt"
	"sync"
	"time"
)

type QueueElement[T any] struct {
	getValueSignal func(types.Unit) (T, error)
	setNextSignal  func(QueueElement[T])
	getNextSignal  func(types.Unit) (QueueElement[T], error)
	j              *junction.Junction
}

func NewQueue[T any]() (func(T), func(types.Unit) (T, error), func()) {
	j := junction.NewJunction()

	firstPort, firstSignal := junction.NewAsyncSignal[QueueElement[T]](j)
	lastPort, lastSignal := junction.NewAsyncSignal[QueueElement[T]](j)
	emptyLastPort, emptyLastSignal := junction.NewSyncSignal[types.Unit, types.Unit](j)

	enqueuePort, enqueueSignal := junction.NewAsyncSignal[T](j)
	dequeuePort, dequeueSignal := junction.NewSyncSignal[types.Unit, T](j)

	tail := newQueueElement[T](QueueElement[T]{})

	junction.NewBinarySyncJoinPattern[QueueElement[T], types.Unit, types.Unit](lastPort, emptyLastPort).Action(func(last QueueElement[T], a types.Unit) types.Unit {
		return types.Unit{}
	})

	junction.NewBinaryAsyncJoinPattern[QueueElement[T], T](lastPort, enqueuePort).Action(func(last QueueElement[T], value T) {
		elem := newQueueElement[T](tail)
		insertIntoElement[T](&elem, value)

		lastSignal(elem)

		if last.getValueSignal != nil {
			last.getNextSignal(types.Unit{})
			last.setNextSignal(elem)
		} else {
			firstSignal(elem)
		}
	})

	junction.NewBinarySyncJoinPattern[QueueElement[T], types.Unit, T](firstPort, dequeuePort).Action(func(first QueueElement[T], a types.Unit) T {
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

	return enqueueSignal, dequeueSignal, func() { junction.Shutdown(j) }
}

func newQueueElement[T any](tail QueueElement[T]) QueueElement[T] {
	j := junction.NewJunction()

	setNextPort, setNextSignal := junction.NewAsyncSignal[QueueElement[T]](j)
	getNextPort, getNextSignal := junction.NewSyncSignal[types.Unit, QueueElement[T]](j)

	junction.NewBinarySyncJoinPattern[types.Unit, QueueElement[T], QueueElement[T]](getNextPort, setNextPort).
		Action(func(a types.Unit, node QueueElement[T]) QueueElement[T] {
			return node
		})

	setNextSignal(tail)

	return QueueElement[T]{
		setNextSignal: setNextSignal,
		getNextSignal: getNextSignal,
		j:             j,
	}
}

func insertIntoElement[T any](elem *QueueElement[T], val T) {
	getValuePort, getValueSignal := junction.NewSyncSignal[types.Unit, T](elem.j)

	junction.NewUnarySyncJoinPattern[types.Unit, T](getValuePort).
		Action(func(a types.Unit) T {
			return val
		})

	(*elem).getValueSignal = getValueSignal
}

func main() {
	enqueue, dequeue, _ := NewQueue[int]()
	producerCount := 1000000
	consumerCount := 900

	var wg sync.WaitGroup
	start := time.Now()

	// Producer
	for i := 1; i <= producerCount; i++ {
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
			//fmt.Println("Consumer", num, " consuming ", val)
		}(i)
	}

	wg.Wait()
	end := time.Since(start)
	fmt.Println("Duration: ", end)
}
