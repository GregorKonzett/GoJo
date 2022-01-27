package unary

import (
	"../../helper"
	"../../types"
	"errors"
)

type AsyncPartialPattern[T any] struct {
	JunctionId int
	Signals    []types.Port
}

type SyncPartialPattern[T any, R any] struct {
	JunctionId int
	Signals    []types.Port
}

func (pattern AsyncPartialPattern[T]) Action(do func(T)) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}

	pattern.Signals[0].JunctionChannel <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Ports:  pattern.Signals,
			Action: helper.WrapUnaryAsync[T](do),
		}},
	}

	return nil
}

func (pattern SyncPartialPattern[T, R]) Action(do func(T) R) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}

	pattern.Signals[0].JunctionChannel <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Ports:  pattern.Signals,
			Action: helper.WrapUnarySync[T, R](do),
		}},
	}

	return nil
}
