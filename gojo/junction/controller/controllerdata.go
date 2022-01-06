package controller

import (
	"../../types"
	"sync"
)

type JoinPatterns struct {
	portIds             int
	joinPatternId       int
	joinPatterns        map[int]types.JoinPatternPacket
	portsToJoinPatterns map[int][]int
	firedPorts          map[int]types.MessageChannel
	fireMutex           sync.Mutex
}

func setupController() JoinPatterns {
	return JoinPatterns{
		portIds:             0,
		joinPatternId:       0,
		joinPatterns:        make(map[int]types.JoinPatternPacket),
		portsToJoinPatterns: make(map[int][]int),
		firedPorts:          make(map[int]types.MessageChannel),
		fireMutex:           sync.Mutex{},
	}
}
