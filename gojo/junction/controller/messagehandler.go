package controller

import (
	"../../types"
	"fmt"
)

func handleMessage(patterns *JoinPatterns, msg types.Packet) {
	fmt.Println("Incoming message: ", msg.Payload.Msg, msg.SignalId.Id)
	(*patterns).firedPorts[msg.SignalId.Id] = append((*patterns).firedPorts[msg.SignalId.Id], msg.Payload)
	joinPattern := findFireableJoinPattern(patterns, msg.SignalId.Id)

	if joinPattern != -1 {
		fire(patterns, joinPattern)
	}
}

func findFireableJoinPattern(patterns *JoinPatterns, port int) int {
	potentialJoinPatterns := (*patterns).portsToJoinPatterns[port]

	for _, pattern := range potentialJoinPatterns {
		valid := true
		for _, signal := range (*patterns).joinPatterns[pattern].Signals {
			if len((*patterns).firedPorts[signal.Id]) == 0 {
				valid = false
				break
			}
		}

		if valid {
			return pattern
		}
	}

	return -1
}

func fire(patterns *JoinPatterns, foundPattern int) {
	pattern := (*patterns).joinPatterns[foundPattern]

	var params []interface{}
	var syncPorts []chan interface{}

	for _, port := range pattern.Signals {
		params = append(params, (*patterns).firedPorts[port.Id][0].Msg)

		if (*patterns).firedPorts[port.Id][0].Ch != nil {
			syncPorts = append(syncPorts, (*patterns).firedPorts[port.Id][0].Ch)
		}

		(*patterns).firedPorts[port.Id] = append((*patterns).firedPorts[port.Id][:0], (*patterns).firedPorts[port.Id][1:]...)
	}

	switch pattern.DoFunction.(type) {
	case types.UnaryAsync:
		fmt.Println("found unary async")
		go (pattern.DoFunction.(types.UnaryAsync))(params[0])
	case types.UnarySync:
		fmt.Println("found unary sync")
		go func() {
			ret := (pattern.DoFunction.(types.UnarySync))(params[0])

			for _, port := range syncPorts {
				port <- ret
			}
		}()
	case types.BinaryAsync:
		fmt.Println("found binary async")
		go (pattern.DoFunction.(types.BinaryAsync))(params[0], params[1])
	case types.BinarySync:
		fmt.Println("found binary sync")
		go func() {
			ret := (pattern.DoFunction.(types.BinarySync))(params[0], params[1])

			for _, port := range syncPorts {
				port <- ret
			}
		}()

	case types.TernaryAsync:
		fmt.Println("found ternary async")
		go (pattern.DoFunction.(types.TernaryAsync))(params[0], params[1], params[2])
	case types.TernarySync:
		fmt.Println("found ternary sync")
		go func() {
			ret := (pattern.DoFunction.(types.TernarySync))(params[0], params[1], params[2])

			for _, port := range syncPorts {
				port <- ret
			}
		}()
	}
}
