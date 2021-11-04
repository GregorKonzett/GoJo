package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
)

func main() {
	j := junction.NewJunction()

	id1, syncSignal := junction.NewSyncSignal[int, string](j)
	id2, syncSignal1 := junction.NewSyncSignal[types.Unit, string](j)

	fmt.Println("Signal created", id1)
	fmt.Println("Signal created", id2)

	returnVal, _ := syncSignal(1)
	returnVal1, _ := syncSignal1(types.Unit{})

	fmt.Println("Client: ", returnVal)
	fmt.Println("Client: ", returnVal1)
}
