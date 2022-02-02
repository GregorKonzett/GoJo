package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
	"time"
)

// Demonstrate that single threaded controller poses a bottle neck (different join patterns listening on different ports) --> will interfere right now
// Check with example: ports A, B: one binary pattern, two unary in single junction --> compare performance and interference between them
func main() {
	j := junction.NewJunction()

	pA, sA := junction.NewAsyncSignal[types.Unit](j)
	pB, sB := junction.NewAsyncSignal[types.Unit](j)
	pC, sC := junction.NewAsyncSignal[types.Unit](j)

	junction.NewBinaryAsyncJoinPattern[types.Unit, types.Unit](pA, pB).Action(func(a types.Unit, b types.Unit) {
		fmt.Println("AB")
		time.Sleep(time.Millisecond * 100)
		sB(types.Unit{})
	})

	junction.NewBinaryAsyncJoinPattern[types.Unit, types.Unit](pA, pC).Action(func(a types.Unit, b types.Unit) {
		fmt.Println("AC")
		time.Sleep(time.Millisecond * 100)
		sC(types.Unit{})
	})

	sB(types.Unit{})
	sC(types.Unit{})

	for true {
		sA(types.Unit{})
		time.Sleep(time.Millisecond * 100)
	}
}

// two binary patterns ab and ac, loop through it b c a and see if both of them fire
// TODO: seem to stop --> investigate
