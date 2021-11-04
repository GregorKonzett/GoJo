package ternary

import "../../types"
import "../../helper"

type AsyncPartialPattern[T any, R any, S any] struct {
	JunctionId int
	Port       chan types.Packet
	Signals    []types.SignalId
}

type SyncPartialPattern[T any, R any, S any] struct {
	JunctionId int
	Port       chan types.Packet
	Signals    []types.SignalId
}

func (pattern AsyncPartialPattern[T, S, R]) ThenDo(do func(T, S, R)) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{
			Msg: types.JoinPatternPacket{
				Signals:    pattern.Signals,
				DoFunction: helper.WrapTernaryAsync[T, S, R](do),
			},
		},
	}
}

func (pattern SyncPartialPattern[T, S, R]) ThenDo(do func(T, S, R) R) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{
			Msg: types.JoinPatternPacket{
				Signals:    pattern.Signals,
				DoFunction: helper.WrapTernarySync[T, S, R](do),
			},
		},
	}
}
