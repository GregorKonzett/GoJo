package controller

import (
	"../../types"
	"sync"
)

type JoinPatterns struct {
	portIds        int
	joinPatternId  int
	joinPatterns   map[int]types.WrappedJoinPattern
	ports          map[int]chan *types.Payload
	messageBitmask int
	fireMutex      sync.Mutex
}

func setupController() JoinPatterns {
	return JoinPatterns{
		portIds:        0,
		joinPatternId:  0,
		joinPatterns:   make(map[int]types.WrappedJoinPattern),
		ports:          make(map[int]chan *types.Payload),
		messageBitmask: 0,
		fireMutex:      sync.Mutex{},
	}
}
