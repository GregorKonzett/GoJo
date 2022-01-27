package controller

import (
	"../../types"
)

func createNewPort(patterns *JoinPatterns, msg types.Packet) {
	channel := make(chan *types.Payload)
	msg.Payload.Ch <- types.PortCreation{
		Ch:       channel,
		SignalId: (*patterns).portIds,
	}

	go handleIncomingMessages(patterns, channel, (*patterns).portIds)

	(*patterns).portIds++
}

/*
	Sends all payloads to each relevant join pattern
*/
func handleIncomingMessages(patterns *JoinPatterns, ch chan *types.Payload, portId int) {
	for true {
		data := <-ch
		(*patterns).portMutex.RLock()

		joinPatterns := (*patterns).portsToJoinPattern[portId]

		for _, pattern := range joinPatterns {
			pattern <- types.WrappedPayload{
				Payload: data,
				PortId:  portId,
			}
		}

		(*patterns).portMutex.RUnlock()
	}
}
