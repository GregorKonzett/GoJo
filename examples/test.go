package main

import (
	"../gojo"
	"fmt"
)

func main() {
	j := gojo.NewJunction()

	id1, syncSignal := gojo.NewSyncSignal[int, string](j)
	id2, syncSignal1 := gojo.NewSyncSignal[int, string](j)

	fmt.Println("Signal created", id1)
	fmt.Println("Signal created", id2)

	returnVal, _ := syncSignal(1)
	returnVal1, _ := syncSignal1(2)

	fmt.Println("Client: ", returnVal)
	fmt.Println("Client: ", returnVal1)

	/*
		Define join patterns:
			j.when(id1).and(id2).thenDo(func(a,b) {})
	*/
}
