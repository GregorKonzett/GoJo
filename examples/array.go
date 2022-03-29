package main

import (
	"fmt"
	"github.com/junctional/GoJo/gojo/junction"
	"github.com/junctional/GoJo/gojo/types"
	"time"
)

type ArrayElement[T any] struct {
	get func(types.Unit) (T, error)
	set func(T)
}

func newArray[T any](size int) []ArrayElement[T] {

	var arrayElems []ArrayElement[T]

	for i := 0; i < size; i++ {
		j := junction.NewJunction()

		getPort, get := junction.NewSyncPort[types.Unit, T](j)
		setPort, set := junction.NewAsyncPort[T](j)

		junction.NewBinarySyncJoinPattern[T, types.Unit, T](setPort, getPort).Action(func(val T, b types.Unit) T {
			return val
		})

		arrayElems = append(arrayElems, ArrayElement[T]{
			get: get,
			set: set,
		})
	}

	return arrayElems
}

func main() {
	size := 5
	array := newArray[int](size)

	// Fill array
	for i := 0; i < size; i++ {
		go func(num int) {
			for true {
				time.Sleep(time.Second)
				array[num].set(num)
				fmt.Println(num, " Producing ", num)
			}

		}(i)
	}

	// Consume array
	for i := 0; i < size; i++ {
		go func(num int) {
			for true {
				time.Sleep(time.Second)
				val, _ := array[num].get(types.Unit{})
				fmt.Println(num, " Consuming ", val)
			}
		}(i)
	}

	for true {
	}
}
