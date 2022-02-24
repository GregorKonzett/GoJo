package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
)

func main() {
	j := junction.NewJunction()

	get, getSignal := junction.NewSyncPort[types.Unit, int](j)
	put, putSignal := junction.NewAsyncPort[int](j)
	full, fullSignal := junction.NewAsyncPort[int](j)
	empty, emptySignal := junction.NewAsyncPort[types.Unit](j)

	junction.NewBinarySyncJoinPattern[types.Unit, int, int](get, full).Action(func(a types.Unit, b int) int {
		fmt.Println("Get body called with ", a, b)
		emptySignal(types.Unit{})
		return b
	})

	junction.NewBinaryAsyncJoinPattern[int, types.Unit](put, empty).Action(func(a int, b types.Unit) {
		fmt.Println("put body called with ", a, b)
		fullSignal(a)
	})

	emptySignal(types.Unit{})
	go putSignal(1234)
	fmt.Println(getSignal(types.Unit{}))

	for true {
	}
}
