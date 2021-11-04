package main

import "fmt"
import "../gojo/junction"

func main() {
	j := junction.NewJunction()

	id1, asyncSignal1 := junction.NewAsyncSignal[int](j)
	id2, asyncSignal2 := junction.NewAsyncSignal[string](j)

	fmt.Println("Signal created", id1)

	junction.NewBinaryAsyncJoinPattern[int, string](j, id1, id2).ThenDo(func(a int, b string) {
		fmt.Println("We got a ", a, " and ", b)
	})

	asyncSignal1(1)
	asyncSignal2("hello there")

	for true {
	}
}
