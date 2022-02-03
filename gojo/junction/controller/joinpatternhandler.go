package controller

import (
	"../../types"
	"sort"
	"time"
)

type Slice struct {
	sort.IntSlice
	idx []int
}

func (s Slice) Swap(i, j int) {
	s.IntSlice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(n []int) *Slice {
	s := &Slice{IntSlice: sort.IntSlice(n), idx: make([]int, len(n))}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

func registerNewJoinPattern(patterns *JoinPatterns, pattern types.JoinPatternPacket) {
	channel := registerJoinPatternWithPorts(patterns, pattern)
	(*patterns).joinPatternId++

	portOrder := make([]int, len(pattern.Ports))

	for i, portId := range pattern.Ports {
		portOrder[i] = portId.Id
	}

	go processJoinPattern(pattern.Action, len(pattern.Ports), channel, portOrder)
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

func processJoinPattern(action interface{}, paramAmount int, ch chan types.WrappedPayload, portOrders []int) {
	allParams := make(map[int][]*types.WrappedPayload, paramAmount)

	s := NewSlice(portOrders)
	sort.Sort(s)

	for true {
		incomingMessage := <-ch

		if _, found := allParams[incomingMessage.PortId]; !found {
			allParams[incomingMessage.PortId] = []*types.WrappedPayload{&incomingMessage}
		} else {
			allParams[incomingMessage.PortId] = append(allParams[incomingMessage.PortId], &incomingMessage)
		}

		params, syncPorts, found := tryClaimMessages(allParams, s.IntSlice, s.idx)

		if found {
			fire(action, params, syncPorts)
		}
	}
}

func fire(action interface{}, params []interface{}, syncPorts []chan interface{}) {
	time.Sleep(time.Millisecond * 10)
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
