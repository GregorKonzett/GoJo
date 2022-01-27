package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
)

func createUnaryPattern[T any](num int, pA types.Port) {
	junction.NewUnaryAsyncJoinPattern[T](pA).Action(func(a T) {
		fmt.Println(num, "Unary")
	})
}

func createBinaryPattern[T any, R any](num int, pA types.Port, pB types.Port) {
	junction.NewBinaryAsyncJoinPattern[T, R](pA, pB).Action(func(a T, b R) {
		fmt.Println(num, "Binary")
	})
}

func createTernaryPattern[T any, R any, S any](num int, pA types.Port, pB types.Port, pC types.Port) {
	junction.NewTernaryAsyncJoinPattern[T, R, S](pA, pB, pC).Action(func(a T, b R, c S) {
		fmt.Println(num, "Ternary")
	})
}

func main() {
	j := junction.NewJunction()

	patternCount := 10
	signalSendingCount := 1000

	var signals []func(types.Unit)

	for i := 0; i < patternCount; i++ {
		pA, sA := junction.NewAsyncSignal[types.Unit](j)
		pB, sB := junction.NewAsyncSignal[types.Unit](j)
		pC, sC := junction.NewAsyncSignal[types.Unit](j)

		signals = append(signals, sA, sB, sC)

		createUnaryPattern[types.Unit](i, pA)
		createUnaryPattern[types.Unit](i, pB)
		createUnaryPattern[types.Unit](i, pC)

		createBinaryPattern[types.Unit, types.Unit](i, pA, pB)
		createBinaryPattern[types.Unit, types.Unit](i, pB, pC)
		createBinaryPattern[types.Unit, types.Unit](i, pA, pC)

		createTernaryPattern[types.Unit, types.Unit, types.Unit](i, pA, pB, pC)
	}

	for i := 0; i < signalSendingCount; i++ {
		for _, s := range signals {
			s(types.Unit{})
		}
	}

	for true {
	}
}
