package main

import (
	"../gojo"
	"fmt"
)

func sendSignal[T any](val T, putSignal func(T)) {
	putSignal(val)
}

func main() {
	j := gojo.NewJunction()

	get, getSignal := gojo.NewSyncSignal[int](j)
	put, putSignal := gojo.NewAsyncSignal[int](j)
	full, fullSignal := gojo.NewAsyncSignal[int](j)
	empty, emptySignal := gojo.NewAsyncSignal[bool](j)

	gojo.NewBinaryRecvJoinPattern[int, int](j, get, full).ThenDo(func(a int) int {
		emptySignal(true)
		return a
	})

	gojo.NewBinarySendJoinPattern[int, interface{}](j, put, empty).ThenDo(func(a int, b interface{}) {
		fullSignal(a)
	})

	go sendSignal(1, putSignal)
	fmt.Println(getSignal())
}
