package ternary

import "../../types"
import "../../helper"

type AsyncPartialPattern[T any, R any, S any] struct {
	JunctionId   int
	Port         chan types.Packet
	InputSignals []types.SignalId
}

type SyncPartialPattern[T any, R any, S any] struct {
	JunctionId    int
	Port          chan types.Packet
	InputSignals  []types.SignalId
	OutputSignals []types.SignalId
}

func (pattern AsyncPartialPattern[T, S, R]) ThenDo(do func(T, S, R)) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Msg: types.JoinPatternPacket{
			InputPorts:  pattern.InputSignals,
			OutputPorts: []types.SignalId{},
			DoFunction:  helper.WrapTernaryAsync[T, S, R](do),
		},
	}
}

func (pattern SyncPartialPattern[T, S, R]) ThenDo(do func(T, S) R) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Msg: types.JoinPatternPacket{
			InputPorts:  pattern.InputSignals,
			OutputPorts: pattern.OutputSignals,
			DoFunction:  helper.WrapTernarySync[T, S, R](do),
		},
	}
}
