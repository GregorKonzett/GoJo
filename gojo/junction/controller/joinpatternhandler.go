package controller

import (
	"../../types"
	"math"
	"reflect"
)

func registerNewJoinPattern(patterns *JoinPatterns, pattern types.JoinPatternPacket) {
	(*patterns).joinPatternId++
	recvChannels := getPortChannels(patterns, pattern)

	go handleIncomingMessages(pattern.Action, recvChannels, (*patterns).joinPatternId)
}

func getPortChannels(patterns *JoinPatterns, joinPattern types.JoinPatternPacket) []chan types.Payload {
	var recvChannels []chan types.Payload

	for _, signal := range joinPattern.Signals {
		recvChannels = append(recvChannels, patterns.ports[signal.Id])
	}

	return recvChannels
}

// TODO: Currently select reads same messages multiple times --> make unique or make queue
func handleIncomingMessages(action interface{}, recvChannels []chan types.Payload, id int) {
	allParams := make([][]types.Payload, len(recvChannels))
	foundAll := 0
	expectedPattern := int(math.Pow(2, float64(len(recvChannels)))) - 1
	cases := make([]reflect.SelectCase, len(recvChannels))

	for i, ch := range recvChannels {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	for true {
		for foundAll&expectedPattern != expectedPattern {
			chosen, value, _ := reflect.Select(cases)

			payload := value.Interface().(types.Payload)

			foundAll |= 1 << chosen

			allParams[chosen] = append(allParams[chosen], payload)
		}

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
