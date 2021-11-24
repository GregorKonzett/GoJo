package queue

import (
	"../../gojo/types"
	"sync"
)

func NewQueueMutex[T any]() (func(T), func(types.Unit) (T, error)) {
	m := sync.Mutex{}
	c := sync.NewCond(&m)

	var arr []T

	enqueue := func(val T) {
		c.L.Lock()
		arr = append(arr, val)
		c.Broadcast()
		c.L.Unlock()
	}

	dequeue := func(types.Unit) (T, error) {
		c.L.Lock()
		for len(arr) == 0 {
			c.Wait()
		}

		val := arr[0]
		arr = arr[1:]

		c.L.Unlock()

		return val, nil
	}

	return enqueue, dequeue
}
