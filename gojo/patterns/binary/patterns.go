package binary

import "../../types"
import "../../helper"

type AsyncPartialPattern[T any, R any] struct {
	JunctionId   int
	Port         chan types.Packet
	InputSignals []types.SignalId
}

type SyncPartialPattern[T any, R any] struct {
	JunctionId    int
	Port          chan types.Packet
	InputSignals  []types.SignalId
	OutputSignals []types.SignalId
}

func (pattern AsyncPartialPattern[T, R]) ThenDo(do func(T, R)) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Msg: types.JoinPatternPacket{
			InputPorts:  pattern.InputSignals,
			OutputPorts: []types.SignalId{},
			DoFunction:  helper.WrapBinaryAsync[T](do),
		},
	}
}

func (pattern SyncPartialPattern[T, R]) ThenDo(do func(T) R) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Msg: types.JoinPatternPacket{
			InputPorts:  pattern.InputSignals,
			OutputPorts: pattern.OutputSignals,
			DoFunction:  helper.WrapBinarySync[T](do),
		},
	}
}
