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
	setValueSignal func(T)
	setNextSignal  func(QueueElement[T])
	getNextSignal  func(types.Unit) (QueueElement[T], error)
	j              *junction.Junction
}

/*
Problem when reading and writing async:
	line 41 will block indefinitely, if first == last and the consumer calls getNextSignal first

	Potential (less ideal) fix: introduce mutex for get getNextSignal fct
*/

func NewQueue[T any]() (func(T), func(types.Unit) (T, error)) {
	j := junction.NewJunction()
	j1 := junction.NewJunction()

	firstPort, firstSignal := junction.NewAsyncSignal[QueueElement[T]](j)
	lastPort, lastSignal := junction.NewAsyncSignal[QueueElement[T]](j1)

	enqueuePort, enqueueSignal := junction.NewAsyncSignal[T](j1)
	dequeuePort, dequeueSignal := junction.NewSyncSignal[types.Unit, T](j)

	junction.NewBinaryAsyncJoinPattern[QueueElement[T], T](lastPort, enqueuePort).Action(func(last QueueElement[T], value T) {
		elem := newQueueElement[T]()

		last.setValueSignal(value)
		last.setNextSignal(elem)

		lastSignal(elem)
	})

	junction.NewBinarySyncJoinPattern[QueueElement[T], types.Unit, T](firstPort, dequeuePort).Action(func(first QueueElement[T], a types.Unit) T {
		nextSignal, _ := first.getNextSignal(types.Unit{})
		firstSignal(nextSignal)

		val, _ := first.getValueSignal(types.Unit{})

		return val
	})

	elem := newQueueElement[T]()

	firstSignal(elem)
	lastSignal(elem)

	return enqueueSignal, dequeueSignal
}

func newQueueElement[T any]() QueueElement[T] {
	j := junction.NewJunction()

	setNextPort, setNextSignal := junction.NewAsyncSignal[QueueElement[T]](j)
	getNextPort, getNextSignal := junction.NewSyncSignal[types.Unit, QueueElement[T]](j)
	getValuePort, getValueSignal := junction.NewSyncSignal[types.Unit, T](j)
	setValuePort, setValueSignal := junction.NewAsyncSignal[T](j)

	junction.NewBinarySyncJoinPattern[types.Unit, QueueElement[T], QueueElement[T]](getNextPort, setNextPort).
		Action(func(a types.Unit, node QueueElement[T]) QueueElement[T] {
			return node
		})

	junction.NewBinarySyncJoinPattern[T, types.Unit, T](setValuePort, getValuePort).
		Action(func(val T, a types.Unit) T {
			// Added for benchmarking
			time.Sleep(time.Duration(10))
			return val
		})

	return QueueElement[T]{
		setNextSignal:  setNextSignal,
		getNextSignal:  getNextSignal,
		setValueSignal: setValueSignal,
		getValueSignal: getValueSignal,
	}
}

func main() {
	enqueue, dequeue := NewQueue[int]()
	producerCount := 2
	consumerCount := 2

	var wg sync.WaitGroup
	start := time.Now()

	// Producer
	for i := 1; i <= producerCount; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()

			for true {
				time.Sleep(750)
				fmt.Println("Producer", num, " Enqueueing ", num)
				enqueue(num)
			}
		}(i)
	}

	// Consumer
	for i := 0; i < consumerCount; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()

			for true {
				time.Sleep(500)
				val, _ := dequeue(types.Unit{})
				fmt.Println("Consumer", num, " consuming ", val)
			}
		}(i)
	}

	wg.Wait()
	end := time.Since(start)
	fmt.Println("Duration: ", end)
}
