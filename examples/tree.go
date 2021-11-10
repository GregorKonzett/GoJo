package main

import (
	"../gojo/junction"
	"../gojo/types"
	"constraints"
)

type Value[T constraints.Ordered, R any] struct {
	key T
	val R
}

type Node[T constraints.Ordered, R any] struct {
	value Value[T, R]
	left  *Node
	right *Node
}

// return left and right signals for next nodes
func insert[T constraints.Ordered, R any](get types.Port, key T, value R) (types.Port, func(T) (R, error), types.Port, func(T) (R, error)) {
	j := junction.NewJunction()
	leftPort, leftSignal := junction.NewSyncSignal[T, R](j)
	rightPort, rightSignal := junction.NewSyncSignal[T, R](j)

	keyPort, keySignal := junction.NewAsyncSignal[T](j)
	valuePort, valueSignal := junction.NewAsyncSignal[R](j)

	keySignal(key)
	valueSignal(value)

	junction.NewTernarySyncJoinPattern[T, R, T, R](keyPort, valuePort, get).Action(func(nodeKey T, nodeValue R, key T) R {
		if nodeKey == key {
			return nodeValue
		} else if nodeKey < key {
			return rightSignal(key)
		} else {
			return leftSignal(key)
		}
	})

	return leftPort, leftSignal, rightPort, rightSignal
}

func newTree[T constraints.Ordered, R any]() (types.Port, func(T) (R, error)) {
	j := junction.NewJunction()

	getPort, getSignal := junction.NewSyncSignal[T, R](j)

	return getPort, getSignal
}

func main() {

}
