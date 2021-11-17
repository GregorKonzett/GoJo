package controller

import "../../types"

func getNewPortId(patterns *JoinPatterns, msg types.Packet) {
	(*patterns).firedPorts[(*patterns).portIds] = []types.Payload{}
	msg.Payload.Ch <- (*patterns).portIds
	(*patterns).portIds++
}
