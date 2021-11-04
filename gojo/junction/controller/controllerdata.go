package controller

import "../../types"

type JoinPatterns struct {
	portIds                int
	joinPatternId          int
	registeredJoinPatterns map[int]types.JoinPatternPacket
	portIdToJoinPatternId  map[int]int
	firedPorts             map[int]int
}

func setupController() JoinPatterns {
	return JoinPatterns{
		portIds:                0,
		joinPatternId:          0,
		registeredJoinPatterns: make(map[int]types.JoinPatternPacket),
		portIdToJoinPatternId:  make(map[int]int),
		firedPorts:             make(map[int]int),
	}
}
