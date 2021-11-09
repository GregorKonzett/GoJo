package ternary

import (
	"../../types"
	"errors"
)
import "../../helper"

type AsyncPartialPattern[T any, S any, R any] struct {
	JunctionId int
	Port       chan types.Packet
	Signals    []types.SignalId
}

type SyncPartialPattern[T any, S any, R any, U any] struct {
	JunctionId int
	Port       chan types.Packet
	Signals    []types.SignalId
}

func (pattern AsyncPartialPattern[T, S, R]) Action(do func(T, S, R)) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}

	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{
			Msg: types.JoinPatternPacket{
				Signals: pattern.Signals,
				Action:  helper.WrapTernaryAsync[T, S, R](do),
			},
		},
	}

	return nil
}

func (pattern SyncPartialPattern[T, S, R, U]) Action(do func(T, S, R) U) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}

	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{
			Msg: types.JoinPatternPacket{
				Signals: pattern.Signals,
				Action:  helper.WrapTernarySync[T, S, R, U](do),
			},
		},
	}

	return nil
}
