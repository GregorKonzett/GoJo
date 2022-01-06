package controller

import (
	"../../types"
	"math/rand"
)

func handleMessage(patterns *JoinPatterns, msg types.Packet) {
	(*patterns).firedPorts[msg.SignalId.Id].Ch <- msg.Payload
	go firePattern(patterns, msg.SignalId.Id)
}

func firePattern(patterns *JoinPatterns, id int) {
	(*patterns).fireMutex.Lock()
	joinPattern := findFireableJoinPattern(patterns, id)

	if joinPattern != -1 {
		fire(patterns, joinPattern)
	}

	(*patterns).fireMutex.Unlock()
}

func findFireableJoinPattern(patterns *JoinPatterns, port int) int {
	potentialJoinPatterns := (*patterns).portsToJoinPatterns[port]

	var validJoinPatterns []int

	for _, pattern := range potentialJoinPatterns {
		valid := true

		for _, signal := range (*patterns).joinPatterns[pattern].Signals {
			if len((*patterns).firedPorts[signal.Id].Ch) == 0 {
				valid = false
				break
			}
		}

		if valid {
			validJoinPatterns = append(validJoinPatterns, pattern)
		}
	}

	if len(validJoinPatterns) > 0 {
		return validJoinPatterns[rand.Intn(len(validJoinPatterns))]
	}

	return -1
}

func fire(patterns *JoinPatterns, foundPattern int) {
	pattern := (*patterns).joinPatterns[foundPattern]

	var params []interface{}
	var syncPorts []chan interface{}

	for _, port := range pattern.Signals {
		param := <-(*patterns).firedPorts[port.Id].Ch
		params = append(params, param.Msg)

		if param.Ch != nil {
			syncPorts = append(syncPorts, param.Ch)
		}
	}

	switch pattern.Action.(type) {
	case types.UnaryAsync:
		go (pattern.Action.(types.UnaryAsync))(params[0])
	case types.UnarySync:
		go func() {
			ret := (pattern.Action.(types.UnarySync))(params[0])

			for _, port := range syncPorts {
				port <- ret
			}
		}()
	case types.BinaryAsync:
		go (pattern.Action.(types.BinaryAsync))(params[0], params[1])
	case types.BinarySync:
		go func() {
			ret := (pattern.Action.(types.BinarySync))(params[0], params[1])

			for _, port := range syncPorts {
				port <- ret
			}
		}()

	case types.TernaryAsync:
		go (pattern.Action.(types.TernaryAsync))(params[0], params[1], params[2])
	case types.TernarySync:
		go func() {
			ret := (pattern.Action.(types.TernarySync))(params[0], params[1], params[2])

			for _, port := range syncPorts {
				port <- ret
			}
		}()
	}
}
