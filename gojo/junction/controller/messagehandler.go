package controller

import (
	"../../types"
)

func handleMessage(patterns *JoinPatterns, msg types.Packet) {
	(*patterns).ports[msg.SignalId.Id] <- msg.Payload
}

/*
func handleMessage1(patterns *JoinPatterns, msg types.Packet) {
	(*patterns).firedPorts[msg.SignalId.Id].Ch <- msg.Payload

	// Set bit corresponding to signal id to 1
	(*patterns).messageBitmask |= 1 << msg.SignalId.Id

	go firePattern(patterns, msg.SignalId.Id)
}

func firePattern(patterns *JoinPatterns, port int) {
	(*patterns).fireMutex.Lock()
	joinPattern := findFireableJoinPattern(patterns, port)

	if joinPattern != -1 {
		fire(patterns, joinPattern)
	}

	(*patterns).fireMutex.Unlock()
}

func findFireableJoinPattern(patterns *JoinPatterns, port int) int {
	potentialJoinPatterns := (*patterns).portsToJoinPatterns[port]
	var validJoinPatterns []int

	for _, pattern := range potentialJoinPatterns {
		if (*patterns).joinPatterns[pattern].Bitmask&(*patterns).messageBitmask == (*patterns).joinPatterns[pattern].Bitmask {
			validJoinPatterns = append(validJoinPatterns, pattern)
		}
	}

	// If last message was consumed in channel, clear bit
	if len(validJoinPatterns) > 0 {
		return validJoinPatterns[rand.Intn(len(validJoinPatterns))]
	}

	return -1
}

func fire(patterns *JoinPatterns, foundPattern int) {
	pattern := (*patterns).joinPatterns[foundPattern]

	var params []interface{}
	var syncPorts []chan interface{}

	for _, port := range pattern.Pattern.Signals {
		param := <-(*patterns).firedPorts[port.Id].Ch
		params = append(params, param.Msg)

		if len((*patterns).firedPorts[port.Id].Ch) == 0 {
			(*patterns).messageBitmask &^= 1 << port.Id
		}

		if param.Ch != nil {
			syncPorts = append(syncPorts, param.Ch)
		}
	}

	switch pattern.Pattern.Action.(type) {
	case types.UnaryAsync:
		go (pattern.Pattern.Action.(types.UnaryAsync))(params[0])
	case types.UnarySync:
		go func() {
			ret := (pattern.Pattern.Action.(types.UnarySync))(params[0])

			for _, port := range syncPorts {
				port <- ret
			}
		}()
	case types.BinaryAsync:
		go (pattern.Pattern.Action.(types.BinaryAsync))(params[0], params[1])
	case types.BinarySync:
		go func() {
			ret := (pattern.Pattern.Action.(types.BinarySync))(params[0], params[1])

			for _, port := range syncPorts {
				port <- ret
			}
		}()

	case types.TernaryAsync:
		go (pattern.Pattern.Action.(types.TernaryAsync))(params[0], params[1], params[2])
	case types.TernarySync:
		go func() {
			ret := (pattern.Pattern.Action.(types.TernarySync))(params[0], params[1], params[2])

			for _, port := range syncPorts {
				port <- ret
			}
		}()
	}
}*/
