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

	junction.NewBinarySyncJoinPattern[int, int](j, get, full).ThenDo(func(a int) int {
		emptySignal(types.Unit{})
		return a
	})

	junction.NewBinaryAsyncJoinPattern[int, types.Unit](j, put, empty).ThenDo(func(a int, b types.Unit) {
		fullSignal(a)
	})

	go putSignal(1)
	fmt.Println(getSignal(types.Unit{}))

	for true {
	}
}
