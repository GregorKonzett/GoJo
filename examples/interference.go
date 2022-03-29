package main

import (
	"fmt"
	"github.com/junctional/GoJo/gojo/junction"
	"github.com/junctional/GoJo/gojo/types"
	"time"
)

// Demonstrate that single threaded controller poses a bottle neck (different join patterns listening on different ports) --> will interfere right now
// Check with example: ports A, B: one binary pattern, two unary in single junction --> compare performance and interference between them
func main() {
	j := junction.NewJunction()

	pA, sA := junction.NewAsyncPort[types.Unit](j)
	pB, sB := junction.NewAsyncPort[types.Unit](j)

	junction.NewBinaryAsyncJoinPattern[types.Unit, types.Unit](pA, pB).Action(func(a types.Unit, b types.Unit) {
		fmt.Println("Binary")
		time.Sleep(time.Millisecond * 100)
		sA(types.Unit{})
		sB(types.Unit{})
	})

	junction.NewUnaryAsyncJoinPattern[types.Unit](pA).Action(func(a types.Unit) {
		fmt.Println("Unary A")
		time.Sleep(time.Millisecond * 100)
		sA(types.Unit{})
	})

	junction.NewUnaryAsyncJoinPattern[types.Unit](pB).Action(func(a types.Unit) {
		fmt.Println("Unary B")
		time.Sleep(time.Millisecond * 100)
		sB(types.Unit{})
	})

	sA(types.Unit{})
	sB(types.Unit{})

	for true {
	}
}
