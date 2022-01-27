package controller

import (
	"../../types"
)

func registerNewJoinPattern(patterns *JoinPatterns, pattern types.JoinPatternPacket) {
	channel := registerJoinPatternWithPorts(patterns, pattern)
	(*patterns).joinPatternId++

	go processJoinPattern(pattern.Action, len(pattern.Ports), channel, pattern.Ports)
}

func registerJoinPatternWithPorts(patterns *JoinPatterns, pattern types.JoinPatternPacket) chan types.WrappedPayload {
	channel := make(chan types.WrappedPayload)

	(*patterns).portMutex.Lock()

	for _, port := range pattern.Ports {
		patternAlreadyIncluded := false
		for _, includedChannel := range (*patterns).portsToJoinPattern[port.Id] {
			if includedChannel == channel {
				patternAlreadyIncluded = true
				break
			}
		}

		if !patternAlreadyIncluded {
			(*patterns).portsToJoinPattern[port.Id] = append((*patterns).portsToJoinPattern[port.Id], channel)
		}
	}

	(*patterns).portMutex.Unlock()

	return channel
}

func processJoinPattern(action interface{}, paramAmount int, ch chan types.WrappedPayload, portOrders []types.Port) {
	allParams := make(map[int][]*types.WrappedPayload, paramAmount)

	for true {
		incomingMessage := <-ch

		if _, found := allParams[incomingMessage.PortId]; !found {
			allParams[incomingMessage.PortId] = []*types.WrappedPayload{&incomingMessage}
		} else {
			allParams[incomingMessage.PortId] = append(allParams[incomingMessage.PortId], &incomingMessage)
		}

		params, syncPorts, found := tryClaimMessages(allParams, portOrders)

		if found {
			fire(action, params, syncPorts)
		}
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
