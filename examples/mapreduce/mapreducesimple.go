package main

import (
	"constraints"
	"fmt"
	"github.com/junctional/GoJo/gojo/junction"
	"github.com/junctional/GoJo/gojo/types"
	"time"
)

func createMapReduce[T constraints.Integer](vals []T) func(types.Unit) ([]T, error) {
	j := junction.NewJunction()

	mapPort, mapSignal := junction.NewAsyncPort[T](j)
	reducePort, reduceSignal := junction.NewAsyncPort[T](j)
	listPort, listSignal := junction.NewAsyncPort[[]T](j)
	finalPort, finalSignal := junction.NewAsyncPort[[]T](j)
	getPort, getSignal := junction.NewSyncPort[types.Unit, []T](j)

	// map
	junction.NewUnaryAsyncJoinPattern[T](mapPort).Action(func(val T) {
		fmt.Println("Mapping: ", val)
		reduceSignal(val * 2)
	})

	// reduce
	junction.NewBinaryAsyncJoinPattern[[]T, T](listPort, reducePort).Action(func(list []T, val T) {
		fmt.Printf("Reducing: %v %d \n", list, val)
		list = append(list, val)

		if len(list) == len(vals) {
			finalSignal(list)
		} else {
			listSignal(list)
		}
	})

	// final
	junction.NewBinarySyncJoinPattern[[]T, types.Unit, []T](finalPort, getPort).Action(func(list []T, a types.Unit) []T {
		return list
	})

	listSignal(make([]T, 0))

	for _, val := range vals {
		go mapSignal(val)
	}

	return getSignal
}

func main() {
	var arr []int
	items := 1000

	for item := items; item > 0; item-- {
		arr = append(arr, item)
	}

	start := time.Now()
	get := createMapReduce[int](arr)
	get(types.Unit{})
	end := time.Since(start)
	fmt.Println("Duration for", items, ": ", end)
}
