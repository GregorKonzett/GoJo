package gojo

import (
	"fmt"
)

func StartController[T any](sender chan Packet[T]) {
	go runThread(sender)
}

func runThread[T any](sender chan Packet[T]) {
	fmt.Println("in thread")
	for true {
		data := <-sender

		fmt.Println(data)
	}
}
