package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
)

func main() {
	j := junction.NewJunction()

	get, getSignal := junction.NewSyncSignal[types.Unit, int](j)
	put, putSignal := junction.NewAsyncSignal[int](j)
	full, fullSignal := junction.NewAsyncSignal[int](j)
	empty, emptySignal := junction.NewAsyncSignal[types.Unit](j)

	junction.NewBinarySyncJoinPattern[types.Unit, int, int](j, get, full).Action(func(a types.Unit, b int) int {
		fmt.Println("Get body called with ", a, b)
		emptySignal(types.Unit{})
		return b
	})

	junction.NewBinaryAsyncJoinPattern[int, types.Unit](j, put, empty).Action(func(a int, b types.Unit) {
		fmt.Println("put body called with ", a, b)
		fullSignal(a)
	})

	emptySignal(types.Unit{})
	go putSignal(1234)
	fmt.Println(getSignal(types.Unit{}))

	for true {
	}
}
