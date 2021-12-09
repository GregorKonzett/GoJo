package main

import (
	".."
	"../../../gojo/types"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func test(producerCount int, consumerCount int, vals int, enqueue func(int), dequeue func(types.Unit) (int, error)) {
	var wg sync.WaitGroup
	start := time.Now()

	// Producer
	for i := 0; i < 10; i++ {
		enqueue(i)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := 1; i <= producerCount; i++ {
			wg.Add(1)
			go func(num int) {
				defer wg.Done()

				for j := 1; j <= vals; j++ {
					enqueue(j * num)
				}

			}(i)
		}
	}()

	wg.Add(1)

	// Consumer
	go func() {
		defer wg.Done()

		for i := 0; i < consumerCount; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for j := 0; j < vals; j++ {
					dequeue(types.Unit{})
				}

			}()
		}
	}()

	wg.Wait()
	end := time.Since(start)
	fmt.Println("Duration: ", end)
}

func main() {
	producerAmount, _ := strconv.Atoi(os.Args[2])
	consumerAmount, _ := strconv.Atoi(os.Args[3])
	vals, _ := strconv.Atoi(os.Args[4])

	if os.Args[1] == "mutex" {
		enqueue, dequeue := queue.NewQueueMutex[int]()
		test(producerAmount, consumerAmount, vals, enqueue, dequeue)
	} else {
		enqueue, dequeue := queue.NewQueue[int]()
		test(producerAmount, consumerAmount, vals, enqueue, dequeue)
	}
}
