package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
	"time"
)

func getReaderWriter[T any]() (func(T), func(types.Unit) (T, error)) {
	j := junction.NewJunction()

	releasePort, produce := junction.NewAsyncPort[T](j)
	acquirePort, consume := junction.NewSyncPort[types.Unit, T](j)

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
		for i := 0; i < 5; i++ {
			time.Sleep(100)
			fmt.Println("Producing: ", val)
			produce(val)
			val += 1
		}
	}()

	// Reader
	for i := 0; i < 3; i++ {
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
