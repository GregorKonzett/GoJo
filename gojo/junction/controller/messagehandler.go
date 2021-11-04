package controller

import (
	"../../types"
	"fmt"
)

func handleMessage(joinPatterns *JoinPatterns, msg types.Packet) {
	fmt.Println("Incoming message: ", msg.Msg)

	if msg.Ch != nil {
		msg.Ch <- "Message processed"
	}
}
