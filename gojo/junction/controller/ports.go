package controller

import "../../types"

func createNewPort(patterns *JoinPatterns, msg types.Packet) {
	channel := make(chan *types.Payload)
	(*patterns).ports[(*patterns).portIds] = channel
	msg.Payload.Ch <- types.PortCreation{
		Ch:       channel,
		SignalId: (*patterns).portIds,
	}
	(*patterns).portIds++
}
