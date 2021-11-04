package controller

import (
	"../../types"
)

func StartController(receiver chan types.Packet) {
	go runThread(receiver)
}

func runThread(receiver chan types.Packet) {
	patterns := setupController()

	for true {
		data := <-receiver
		switch data.Type {
		case types.MESSAGE:
			handleMessage(&patterns, data)
		case types.AddJoinPattern:
			registerNewJoinPattern(&patterns, data.Payload.Msg.(types.JoinPatternPacket))
		case types.GetNewPortId:
			getNewPortId(&patterns, data)
		}
	}
}
