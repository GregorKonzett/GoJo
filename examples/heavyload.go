package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
)

func createUnaryPattern[T any](num int, pA types.Port) {
	junction.NewUnaryAsyncJoinPattern[T](pA).Action(func(a T) {
		fmt.Println(num, "Unary", a)
	})
}

func createBinaryPattern[T any, R any](num int, pA types.Port, pB types.Port) {
	junction.NewBinaryAsyncJoinPattern[T, R](pA, pB).Action(func(a T, b R) {
		fmt.Println(num, "Binary", a, b)
	})
}

func createTernaryPattern[T any, R any, S any](num int, pA types.Port, pB types.Port, pC types.Port) {
	junction.NewTernaryAsyncJoinPattern[T, R, S](pA, pB, pC).Action(func(a T, b R, c S) {
		fmt.Println(num, "Ternary", a, b, c)
	})
}

func main() {
	j := junction.NewJunction()

	patternCount := 2
	signalSendingCount := 2

	var signals []func(int)

	for i := 0; i < patternCount; i++ {
		pA, sA := junction.NewAsyncPort[int](j)
		pB, sB := junction.NewAsyncPort[int](j)
		pC, sC := junction.NewAsyncPort[int](j)

		signals = append(signals, sA, sB, sC)

		createUnaryPattern[int](i*10+1, pA)
		createUnaryPattern[int](i*10+2, pB)
		createUnaryPattern[int](i*10+3, pC)

		createBinaryPattern[int, int](i, pA, pB)
		createBinaryPattern[int, int](i, pB, pC)
		createBinaryPattern[int, int](i, pA, pC)

		createTernaryPattern[int, int, int](i, pA, pB, pC)
	}

	for i := 0; i < signalSendingCount; i++ {
		for j, s := range signals {
			s(j)
		}
	}

	for true {
	}
}
