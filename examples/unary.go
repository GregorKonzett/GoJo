package main

import "fmt"
import "../gojo/junction"

func main() {
	j := junction.NewJunction()

	id1, asyncSignal := junction.NewAsyncSignal[int](j)

	fmt.Println("Signal created", id1)

	junction.NewUnaryAsyncJoinPattern[int](j, id1).ThenDo(func(a int) {
		fmt.Println("We got a ", a)
	})

	asyncSignal(1)
}
