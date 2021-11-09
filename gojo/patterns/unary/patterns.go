package unary

import (
	"../../helper"
	"../../types"
	"errors"
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

func (pattern AsyncPartialPattern[T]) Action(do func(T)) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}

	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Signals: pattern.Signals,
			Action:  helper.WrapUnaryAsync[T](do),
		}},
	}

	return nil
}

func (pattern SyncPartialPattern[T, R]) Action(do func(T) R) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}

	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Signals: pattern.Signals,
			Action:  helper.WrapUnarySync[T, R](do),
		}},
	}

	return nil
}
