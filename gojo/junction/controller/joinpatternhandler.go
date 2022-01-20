package controller

import (
	"../../types"
	"reflect"
)

func registerNewJoinPattern(patterns *JoinPatterns, pattern types.JoinPatternPacket) {
	(*patterns).joinPatternId++
	recvChannels := getPortChannels(patterns, pattern)

	go handleIncomingMessages(pattern.Action, recvChannels)
}

// TODO: Add sync return channels
func getPortChannels(patterns *JoinPatterns, joinPattern types.JoinPatternPacket) []chan types.Payload {
	var recvChannels []chan types.Payload

	for _, signal := range joinPattern.Signals {
		recvChannels = append(recvChannels, patterns.ports[signal.Id])
	}

	return recvChannels
}

// TODO: Currently select reads same messages multiple times --> make unique or make queue
func handleIncomingMessages(action interface{}, recvChannels []chan types.Payload) {
	for true {
		cases := make([]reflect.SelectCase, len(recvChannels))

		for i, ch := range recvChannels {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}

		remaining := len(cases)
		params := make([]interface{}, len(recvChannels))
		var syncPorts []chan interface{}

		for remaining > 0 {
			chosen, value, _ := reflect.Select(cases)

			payload := value.Interface().(types.Payload)

			params[chosen] = payload.Msg

			if payload.Ch != nil {
				syncPorts = append(syncPorts, payload.Ch)
			}

			remaining--
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
