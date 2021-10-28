package unary

import (
	"../../helper"
	"../../types"
)

type SendPartialPattern[T any] struct {
	JunctionId int
	Port       chan types.Packet
	SignalOne  types.SignalId
}

type SendPattern[T any] struct {
	JunctionId int
	Port       chan types.Packet
	SignalOne  types.SignalId
	Do         func(interface{})
}

type SyncPartialPattern[T any, R any] struct {
	JunctionId int
	SenderOne  types.SignalId
	Receiver   types.SignalId
}

type SyncPattern[T any, R any] struct {
	JunctionId int
	Receiver   types.SignalId
	Do         func(T) R
}

func (pattern SendPartialPattern[T]) ThenDo(do func(T)) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Msg: SendPattern[T]{
			JunctionId: pattern.JunctionId,
			Port:       pattern.Port,
			SignalOne:  pattern.SignalOne,
			Do:         helper.WrapUnarySend[T](do),
		},
	}
}
