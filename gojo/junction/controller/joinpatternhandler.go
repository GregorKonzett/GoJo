package controller

import (
	"../../types"
)

func registerNewJoinPattern(patterns *JoinPatterns, pattern types.JoinPatternPacket) {
	channel := registerJoinPatternWithPorts(patterns, pattern)
	(*patterns).joinPatternId++

	portMapping := portToParameterIndex(pattern.Ports)

	go processJoinPattern(pattern.Action, len(pattern.Ports), channel, portMapping)
}

func registerJoinPatternWithPorts(patterns *JoinPatterns, pattern types.JoinPatternPacket) chan types.WrappedPayload {
	channel := make(chan types.WrappedPayload)

	(*patterns).portMutex.Lock()

	for _, port := range pattern.Ports {
		(*patterns).portsToJoinPattern[port.Id] = append((*patterns).portsToJoinPattern[port.Id], channel)
	}

	(*patterns).portMutex.Unlock()

	return channel
}

// TODO: Change this to int array to handle same port multiple times
func portToParameterIndex(ports []types.Port) map[int]int {
	mapping := make(map[int]int)

	for i, port := range ports {
		mapping[port.Id] = i
	}

	return mapping
}

func processJoinPattern(action interface{}, paramAmount int, ch chan types.WrappedPayload, portMapping map[int]int) {
	allParams := make([][]*types.Payload, paramAmount)
	foundAll := 0

	expectedPattern := 1<<paramAmount - 1

	for true {
		// TODO: Check if enough things on each channel (if same channel twice or smth)
		// Mention in report --> new feature to support same channel multiple times
		for foundAll&expectedPattern != expectedPattern {
			incomingMessage := <-ch
			foundAll |= 1 << incomingMessage.PortId

			allParams[portMapping[incomingMessage.PortId]] = append(allParams[portMapping[incomingMessage.PortId]], incomingMessage.Payload)
		}

		// TODO: tryClaim
		var params []interface{}
		var syncPorts []chan interface{}

		for i := 0; i < len(allParams); i++ {
			params = append(params, allParams[i][0].Msg)

			if allParams[i][0].Ch != nil {
				syncPorts = append(syncPorts, allParams[i][0].Ch)
			}

			if len(allParams[i]) == 1 {
				foundAll &^= 1 << i
				allParams[i] = nil
			} else {
				allParams[i] = allParams[i][1:]
			}
		}

		fire(action, params, syncPorts)
	}
}

func fire(action interface{}, params []interface{}, syncPorts []chan interface{}) {
	switch action.(type) {
	case types.UnaryAsync:
		go (action.(types.UnaryAsync))(params[0])
	case types.UnarySync:
		go func() {
			ret := (action.(types.UnarySync))(params[0])

			for _, port := range syncPorts {
				port <- ret
			}
		}()
	case types.BinaryAsync:
		go (action.(types.BinaryAsync))(params[0], params[1])
	case types.BinarySync:
		go func() {
			ret := (action.(types.BinarySync))(params[0], params[1])

			for _, port := range syncPorts {
				port <- ret
			}
		}()

	case types.TernaryAsync:
		go (action.(types.TernaryAsync))(params[0], params[1], params[2])
	case types.TernarySync:
		go func() {
			ret := (action.(types.TernarySync))(params[0], params[1], params[2])

			for _, port := range syncPorts {
				port <- ret
			}
		}()
	}
}

// TODO IDEAS:
/*
	Fan in Fan out messages so we don't have to listen on a dynamic list of channels here
	PROBLEM: No message stealing --> messages are getting lost
*/

// IDEAS:
// sleeps when pattern matching
