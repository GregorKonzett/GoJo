package unary

import (
	"errors"
	"github.com/junctional/GoJo/gojo/helper"
	"github.com/junctional/GoJo/gojo/types"
)

// AsyncPattern Struct containing all ports the join pattern is listening on and defines the Action data types
type AsyncPattern[T any] struct {
	Signals []types.Port
}

// SyncPattern Struct containing all ports the join pattern is listening on and defines the Action data types
type SyncPattern[T any, R any] struct {
	Signals []types.Port
}

// Action Takes a function with the data types defined in the struct that will be executed when the pattern fires
//and registers the pattern with the junction controller
func (pattern AsyncPattern[T]) Action(do func(T)) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}
	resp := make(chan interface{})

	pattern.Signals[0].JunctionChannel <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Ports:  pattern.Signals,
			Action: helper.WrapUnaryAsync[T](do),
		},
			Ch: resp,
		},
	}
	<-resp
	return nil
}

// Action Takes a function with the data types defined in the struct that will be executed when the pattern fires
//and registers the pattern with the junction controller
func (pattern SyncPattern[T, R]) Action(do func(T) R) error {
	if !helper.CheckForSameJunction(pattern.Signals) {
		return errors.New("signals from different junctions")
	}
	resp := make(chan interface{})

	pattern.Signals[0].JunctionChannel <- types.Packet{
		Type: types.AddJoinPattern,
		Payload: types.Payload{Msg: types.JoinPatternPacket{
			Ports:  pattern.Signals,
			Action: helper.WrapUnarySync[T, R](do),
		},
			Ch: resp,
		},
	}
	<-resp
	return nil
}
