package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
)

func getMutex() (func(types.Unit), func(types.Unit) (types.Unit, error)) {
	j := junction.NewJunction()

	releasePort, release := junction.NewAsyncPort[types.Unit](j)
	acquirePort, acquire := junction.NewSyncPort[types.Unit, types.Unit](j)
	lockPort, lock := junction.NewAsyncPort[types.Unit](j)

	junction.NewBinarySyncJoinPattern[types.Unit, types.Unit, types.Unit](lockPort, acquirePort).Action(func(a types.Unit, b types.Unit) types.Unit {
		return types.Unit{}
	})

	junction.NewUnaryAsyncJoinPattern[types.Unit](releasePort).Action(func(a types.Unit) {
		lock(types.Unit{})
	})

	lock(types.Unit{})

	return release, acquire
}

func main() {
	release, acquire := getMutex()

	sharedVar := 0

	for i := 0; i < 10; i += 2 {
		go func() {
			acquire(types.Unit{})
			sharedVar += 2
			fmt.Println("Incrementing: ", sharedVar)
			release(types.Unit{})
		}()
	}

	for i := 1; i < 10; i += 2 {
		go func() {
			acquire(types.Unit{})
			sharedVar -= 2
			fmt.Println("Decrementing: ", sharedVar)
			release(types.Unit{})
		}()
	}

	for true {
	}
}
