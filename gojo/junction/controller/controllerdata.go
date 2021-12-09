package controller

import (
	"../../types"
)

type JoinPatterns struct {
	portIds             int
	joinPatternId       int
	joinPatterns        map[int]types.JoinPatternPacket
	portsToJoinPatterns map[int][]int
	firedPorts          map[int][]types.Payload
}

func setupController() JoinPatterns {
	return JoinPatterns{
		portIds:             0,
		joinPatternId:       0,
		joinPatterns:        make(map[int]types.JoinPatternPacket),
		portsToJoinPatterns: make(map[int][]int),
		firedPorts:          make(map[int][]types.Payload),
	}
}
