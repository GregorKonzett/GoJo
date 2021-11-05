package binary

import "../../types"
import "../../helper"

type AsyncPartialPattern[T any, R any] struct {
	JunctionId int
	Port       chan types.Packet
	Signals    []types.SignalId
}

type SyncPartialPattern[T any, S any, R any] struct {
	JunctionId int
	Port       chan types.Packet
	Signals    []types.SignalId
}

func (pattern AsyncPartialPattern[T, R]) ThenDo(do func(T, R)) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Signals:    pattern.Signals,
			DoFunction: helper.WrapBinaryAsync[T, R](do),
		},
		},
	}
}

func (pattern SyncPartialPattern[T, S, R]) ThenDo(do func(T, S) R) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{
			Msg: types.JoinPatternPacket{
				Signals:    pattern.Signals,
				DoFunction: helper.WrapBinarySync[T, S, R](do),
			},
		},
	}
}
