package main

import (
	"../gojo/junction"
	"../gojo/types"
	"constraints"
	"fmt"
)

func createMapReduce[T constraints.Integer](vals []T) func(types.Unit) ([]T, error) {
	j := junction.NewJunction()

	mapPort, mapSignal := junction.NewAsyncSignal[T](j)
	reducePort, reduceSignal := junction.NewAsyncSignal[T](j)
	listPort, listSignal := junction.NewAsyncSignal[[]T](j)
	finalPort, finalSignal := junction.NewAsyncSignal[[]T](j)
	getPort, getSignal := junction.NewSyncSignal[types.Unit, []T](j)

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
	arr := []int{1, 2, 3, 5}
	get := createMapReduce[int](arr)

	ret, _ := get(types.Unit{})

	fmt.Printf("Result: %v", ret)
}
