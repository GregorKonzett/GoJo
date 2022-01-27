package controller

import (
	"../../types"
	"sync"
)

type JoinPatterns struct {
	portIds            int
	joinPatternId      int
	joinPatterns       map[int]types.WrappedJoinPattern
	portsToJoinPattern map[int][]chan types.WrappedPayload
	portMutex          sync.RWMutex
}

func setupController() JoinPatterns {
	return JoinPatterns{
		portIds:            0,
		joinPatternId:      0,
		joinPatterns:       make(map[int]types.WrappedJoinPattern),
		portsToJoinPattern: make(map[int][]chan types.WrappedPayload),
		portMutex:          sync.RWMutex{},
	}
}
