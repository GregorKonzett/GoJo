package controller

import (
	"../../types"
	"sync"
)

// JoinPatterns will hold the next available portId + joinPatternId and includes a map mapping a portId to all join
// patterns waiting on messages on this port. This Map is safe for concurrent access by using the portMutex mutex.
type JoinPatterns struct {
	portIds            int
	portsToJoinPattern map[int][]chan *types.Packet
	portMutex          sync.RWMutex
}

// setupController initializes the JoinPatterns struct with default values
func setupController() JoinPatterns {
	return JoinPatterns{
		portIds:            0,
		portsToJoinPattern: make(map[int][]chan *types.Packet),
		portMutex:          sync.RWMutex{},
	}
}
