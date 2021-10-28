package binary

import "../../types"
import "../../helper"

type SendPartialPattern[T any, R any] struct {
	JunctionId int
	Port       chan types.Packet
	SignalOne  types.SignalId
	SignalTwo  types.SignalId
}

type SendPattern[T any, R any] struct {
	JunctionId int
	Port       chan types.Packet
	SignalOne  types.SignalId
	SignalTwo  types.SignalId
	Do         func(interface{}, interface{})
}

type RecvPartialPattern[T any, R any] struct {
	JunctionId int
	Port       chan types.Packet
	SignalOne  types.SignalId
	SignalTwo  types.SignalId
}

type RecvPattern[T any, R any] struct {
	JunctionId int
	Port       chan types.Packet
	SignalOne  types.SignalId
	SignalTwo  types.SignalId
	Do         func(interface{}) interface{}
}

func (pattern SendPartialPattern[T, R]) ThenDo(do func(T, R)) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Msg: SendPattern[T, R]{
			JunctionId: pattern.JunctionId,
			Port:       pattern.Port,
			SignalOne:  pattern.SignalOne,
			Do:         helper.WrapBinarySend[T, R](do),
		},
	}
}

func (pattern RecvPartialPattern[T, R]) ThenDo(do func(T) R) {
	pattern.Port <- types.Packet{
		Type: types.AddJoinPattern,
		Msg: RecvPattern[T, R]{
			JunctionId: pattern.JunctionId,
			Port:       pattern.Port,
			SignalOne:  pattern.SignalOne,
			Do:         helper.WrapBinaryRecv[T, R](do),
		},
	}
}
