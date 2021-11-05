package unary

import (
	"../../helper"
	"../../types"
)

type AsyncPartialPattern[T any] struct {
	JunctionId int
	Port       chan types.Packet
	Signals    []types.SignalId
}

type SyncPartialPattern[T any, R any] struct {
	JunctionId int
	Port       chan types.Packet
	Signals    []types.SignalId
}

func (pattern AsyncPartialPattern[T]) ThenDo(do func(T)) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Signals:    pattern.Signals,
			DoFunction: helper.WrapUnaryAsync[T](do),
		}},
	}
}

func (pattern SyncPartialPattern[T, R]) ThenDo(do func(T) R) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Signals:    pattern.Signals,
			DoFunction: helper.WrapUnarySync[T, R](do),
		}},
	}
}
