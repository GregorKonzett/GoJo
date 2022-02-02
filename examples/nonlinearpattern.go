package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
	"time"
)

func main() {
	j := junction.NewJunction()

	pA, sA := junction.NewAsyncSignal[types.Unit](j)
	pB, sB := junction.NewAsyncSignal[types.Unit](j)

	junction.NewTernaryAsyncJoinPattern[types.Unit, types.Unit, types.Unit](pA, pA, pB).Action(func(a types.Unit, b types.Unit, c types.Unit) {
		fmt.Println("Ternary")
		time.Sleep(time.Millisecond * 100)
		sA(types.Unit{})
		sA(types.Unit{})
		sB(types.Unit{})
	})

	sA(types.Unit{})
	sA(types.Unit{})
	sB(types.Unit{})

	for true {
	}
}
