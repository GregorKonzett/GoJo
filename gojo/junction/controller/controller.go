package controller

import (
	"../../types"
)

// StartController starts the controller goroutine, which will register all join patterns + ports for one junction
func StartController(receiver chan types.Packet) {
	go runThread(receiver)
}

// runThread sets up the empty controller data struct and waits for new Packets to register new join patterns or ports.
// It will do so until it receives a Shutdown Packet
func runThread(receiver chan types.Packet) {
	patterns := setupController()

	for true {
		data := <-receiver
		switch data.Type {
		case types.AddJoinPattern:
			registerNewJoinPattern(&patterns, data.Payload.Msg.(types.JoinPatternPacket))
		case types.CreateNewPort:
			createNewPort(&patterns, data)
		case types.Shutdown:
			break
		}
	}
}
