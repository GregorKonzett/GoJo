package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
	"time"
)

func getReaderWriter[T any]() (func(T), func(types.Unit) (T, error)) {
	j := junction.NewJunction()

	releasePort, produce := junction.NewAsyncSignal[T](j)
	acquirePort, consume := junction.NewSyncSignal[types.Unit, T](j)

	junction.NewBinarySyncJoinPattern[T, types.Unit, T](releasePort, acquirePort).Action(func(value T, b types.Unit) T {
		return value
	})

	return produce, consume
}

func main() {
	produce, consume := getReaderWriter[int]()

	// Writer
	go func() {
		val := 0
		for true {
			time.Sleep(100)
			fmt.Println("Producing: ", val)
			produce(val)
			val += 1
		}
	}()

	// Reader
	for i := 1; i < 10; i += 2 {
		go func(num int) {
			for true {
				val, _ := consume(types.Unit{})

				fmt.Println(num, " consuming : ", val)
			}
		}(i)
	}

	for true {
	}
}
