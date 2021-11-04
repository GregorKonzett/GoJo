package unary

import (
	"../../helper"
	"../../types"
)

type AsyncPartialPattern[T any] struct {
	JunctionId   int
	Port         chan types.Packet
	InputSignals []types.SignalId
}

type SyncPartialPattern[R any] struct {
	JunctionId    int
	Port          chan types.Packet
	OutputSignals []types.SignalId
}

func (pattern AsyncPartialPattern[T]) ThenDo(do func(T)) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Msg: types.JoinPatternPacket{
			InputPorts:  pattern.InputSignals,
			OutputPorts: []types.SignalId{},
			DoFunction:  helper.WrapUnaryAsync[T](do),
		},
	}
}

func (pattern SyncPartialPattern[R]) ThenDo(do func() R) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Msg: types.JoinPatternPacket{
			InputPorts:  []types.SignalId{},
			OutputPorts: pattern.OutputSignals,
			DoFunction:  helper.WrapUnarySync[R](do),
		},
	}
}
