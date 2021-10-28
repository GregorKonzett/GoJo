package ternary

import "../../types"

type SendPartialPattern[T any] struct {
	JunctionId  int
	SenderOne   types.SignalId
	SenderTwo   types.SignalId
	SenderThree types.SignalId
}

type SendPattern[T any] struct {
	JunctionId  int
	SenderOne   types.SignalId
	SenderTwo   types.SignalId
	SenderThree types.SignalId
	Do          func(T)
}
