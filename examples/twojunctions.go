package main

import (
	"../gojo/junction"
	"fmt"
	"time"
)

func main() {
	j1 := junction.NewJunction()
	j2 := junction.NewJunction()

	s1, signal1 := junction.NewAsyncSignal[int](j1)
	s2, signal2 := junction.NewAsyncSignal[int](j2)

	junction.NewUnaryAsyncJoinPattern[int](j1, s1).Action(func(a int) {
		fmt.Println("Junction1: ", a)
		time.Sleep(1 * time.Second)
		signal2(a + 1)
	})

	junction.NewUnaryAsyncJoinPattern[int](j2, s2).Action(func(a int) {
		fmt.Println("Junction2: ", a)
		time.Sleep(1 * time.Second)
		signal1(a + 1)
	})

	signal1(0)

	for true {
	}
}
