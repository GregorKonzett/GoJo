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
		case types.Shutdown:
			break
		}
	}
}

// TODO: Check how scalable join pattern research paper handles join patterns with duplicate ports
// Add natural number vector to each join pattern ensuring that there are n messages in the channel in addition to the bitmask
