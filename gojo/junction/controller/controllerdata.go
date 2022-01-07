package controller

import (
	"../../types"
	"sync"
)

type JoinPatterns struct {
	portIds             int
	joinPatternId       int
	joinPatterns        map[int]types.WrappedJoinPattern
	portsToJoinPatterns map[int][]int
	firedPorts          map[int]types.MessageChannel
	messageBitmask      int
	fireMutex           sync.Mutex
}

func setupController() JoinPatterns {
	return JoinPatterns{
		portIds:             0,
		joinPatternId:       0,
		joinPatterns:        make(map[int]types.WrappedJoinPattern),
		portsToJoinPatterns: make(map[int][]int),
		firedPorts:          make(map[int]types.MessageChannel),
		messageBitmask:      0,
		fireMutex:           sync.Mutex{},
	}
}
