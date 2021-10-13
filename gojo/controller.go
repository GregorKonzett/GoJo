package gojo

import (
	"fmt"
)

func StartController(receiver chan interface{}, sender chan interface{}) {
	go runThreadGeneric(receiver, sender)
}

func runThreadGeneric(receiver chan interface{}, sender chan interface{}) {
	for true {
		data := <-receiver

		fmt.Println("Controller: ", data)

		sender <- "Return Value"
	}
}
