package controller

import (
	"../../types"
)

// createNewPort creates a new channel that will be used to receive messages on this port. Additionally, it starts a new
// goroutine in the background waiting for messages arriving on this port
func createNewPort(patterns *JoinPatterns, msg types.Packet) {
	channel := make(chan *types.Payload)
	msg.Payload.Ch <- types.PortCreation{
		Ch:     channel,
		PortId: (*patterns).portIds,
	}

	go handleIncomingMessages(patterns, channel, (*patterns).portIds)

	(*patterns).portIds++
}

// handleIncomingMessages redirects each received Payload to every Join Pattern that is waiting for messages on this
// port
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
