package controller

import (
	"../../types"
)

func handleMessage(patterns *JoinPatterns, msg types.Packet) {
	(*patterns).ports[msg.SignalId.Id] <- msg.Payload
}
