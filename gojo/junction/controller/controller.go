package controller

import (
	"../../types"
	"fmt"
)

func StartController(receiver chan types.Packet) {
	go runThread(receiver)
}

func runThread(receiver chan types.Packet) {
	joinPatterns := setupController()

	for true {
		data := <-receiver
		fmt.Println(joinPatterns)
		switch data.Type {
		case types.MESSAGE:
			handleMessage(&joinPatterns, data)
		case types.AddJoinPattern:
			registerNewJoinPattern(&joinPatterns, data.Msg.(types.JoinPatternPacket))
		case types.GetNewPortId:
			getNewPortId(&joinPatterns, data)
		}
	}
}
