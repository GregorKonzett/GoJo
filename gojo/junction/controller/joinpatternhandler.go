package controller

import (
	"github.com/junctional/GoJo/gojo/types"
)

// registerNewJoinPattern registers the join pattern with each of the ports it's listening on. Additionally, start a new
// goroutine in the background handling all incoming messages to this join pattern and potentially firing it.
func registerNewJoinPattern(patterns *JoinPatterns, pattern types.JoinPatternPacket, ch chan interface{}) {
	channel := registerJoinPatternWithPorts(patterns, pattern)

	portOrder := make([]int, len(pattern.Ports))

	for i, portId := range pattern.Ports {
		portOrder[i] = portId.Id
	}

	ch <- types.Unit{}
	go processJoinPattern(pattern.Action, len(pattern.Ports), channel, portOrder)
}

// registerJoinPatternWithPorts adds the join pattern to each port list and returns a channel that will receive all
// messages sent to the join pattern from all ports
func registerJoinPatternWithPorts(patterns *JoinPatterns, pattern types.JoinPatternPacket) chan *types.Packet {
	channel := make(chan *types.Packet)

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

// processJoinPattern waits for new messages received from any of the ports the join pattern is listening on
// Whenever a new message was received it is appended to a list created for each port and then checked if this join
// pattern can now be fired. If a message can be consumed on each port, the join pattern is fired and it's listening for
// new messages again.
func processJoinPattern(action interface{}, paramAmount int, ch chan *types.Packet, portOrders []int) {
	allParams := make(map[int][]*types.Packet, paramAmount)

	for true {
		incomingMessage := <-ch

		if _, found := allParams[incomingMessage.PortId]; !found {
			allParams[incomingMessage.PortId] = []*types.Packet{incomingMessage}
		} else {
			allParams[incomingMessage.PortId] = append(allParams[incomingMessage.PortId], incomingMessage)
		}

		params, syncPorts, found := tryClaimMessages(allParams, portOrders)

		if found {
			fire(action, params, syncPorts)
		}
	}
}

// fire takes the join pattern's action function, the list of parameters and a list of all syncPorts waiting for a response
// Since the data type of the action isn't known at this point anymore, the arity of the function is first determined
// before executed in it's own goroutine. Once this goroutine completes, the return value is sent to each syncPorts
// (for synchronous join patterns)
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
