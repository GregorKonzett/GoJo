package main

import "fmt"
import "../gojo/junction"

func main() {
	j := junction.NewJunction()

	ch1 := make(chan int)
	ch2 := ch1

	asd := make(map[chan int]int)

	asd[ch1] = 2

	fmt.Println(asd)

	if ch1 == ch2 {
		fmt.Println("equal")
	} else {
		fmt.Println("not equakl")
	}

	id1, asyncSignal1 := junction.NewAsyncSignal[int](j)
	id2, asyncSignal2 := junction.NewAsyncSignal[string](j)

	fmt.Println("Signal created", id1)

	junction.NewBinaryAsyncJoinPattern[int, string](j, id1, id2).Action(func(a int, b string) {
		fmt.Println("We got a ", a, " and ", b)
	})

	asyncSignal1(1)
	asyncSignal2("hello there")

	for true {
	}
}
