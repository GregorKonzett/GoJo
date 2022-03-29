package controller

import (
	"github.com/junctional/GoJo/gojo/types"
)

// createNewPort creates a new channel that will be used to receive messages on this port. Additionally, it starts a new
// goroutine in the background waiting for messages arriving on this port
func createNewPort(patterns *JoinPatterns, msg types.Packet) {
	channel := make(chan types.Packet)
	msg.Payload.Ch <- types.PortCreation{
		Ch:     channel,
		PortId: (*patterns).portIds,
	}

	go handleIncomingSignals(patterns, channel)

	(*patterns).portIds++
}

// handleIncomingSignals redirects each received Payload to every Join Pattern that is waiting for messages on this
// port
func handleIncomingSignals(patterns *JoinPatterns, ch chan types.Packet) {
	for true {
		data := <-ch
		(*patterns).portMutex.RLock()

		joinPatterns := (*patterns).portsToJoinPattern[data.PortId]

		for _, pattern := range joinPatterns {
			pattern <- &data
		}

		(*patterns).portMutex.RUnlock()
	}
}
