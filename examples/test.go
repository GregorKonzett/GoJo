package main

import (
	"../gojo"
)

func main() {
	sender := make(chan gojo.Packet[any])

	gojo.StartController(sender)

	sender <- gojo.Packet[any]{
		Msg: gojo.Message[any]{
			Data: 1,
		},
	}

	sender <- gojo.Packet[any]{
		Msg: gojo.Message[any]{
			Data: "asd",
		},
	}
}
