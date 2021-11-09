package binary

import (
	"../../types"
	"errors"
)
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

func (pattern AsyncPartialPattern[T, R]) Action(do func(T, R)) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}

	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Signals: pattern.Signals,
			Action:  helper.WrapBinaryAsync[T, R](do),
		},
		},
	}

	return nil
}

func (pattern SyncPartialPattern[T, S, R]) Action(do func(T, S) R) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}

	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{
			Msg: types.JoinPatternPacket{
				Signals: pattern.Signals,
				Action:  helper.WrapBinarySync[T, S, R](do),
			},
		},
	}

	return nil
}
