package main

import (
	".."
	"../../../gojo/types"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func test(producerCount int, consumerCount int, enqueue func(int), dequeue func(types.Unit) (int, error)) {
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

func main() {
	amount := 10
	producerAmount := 100000
	consumerAmount := 1000
	fmt.Println("Multiple Junctions:")

	for i := 0; i < amount; i++ {
		fmt.Println("Active goroutines: ", runtime.NumGoroutine())
		enqueue, dequeue, shutdown := queue.NewQueue[int]()
		test(producerAmount, consumerAmount, enqueue, dequeue)
		shutdown()
		time.Sleep(time.Second * 10)
	}

	fmt.Println("1 Junction:")

	for i := 0; i < amount; i++ {
		fmt.Println("Active goroutines: ", runtime.NumGoroutine())
		enqueue, dequeue, shutdown := queue.NewQueue1[int]()
		test(producerAmount, consumerAmount, enqueue, dequeue)
		shutdown()
		time.Sleep(time.Second * 10)
	}
}
