package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
	"math/rand"
	"time"
)

type Fork struct {
	Id   types.SignalId
	Free func(types.Unit)
}

type Philosopher struct {
	Id  types.SignalId
	Eat func(types.Unit)
}

func main() {
	philosoperCount := 5
	j := junction.NewJunction()

	var forks []Fork
	var philosophers []Philosopher

	// Create Forks
	for i := 0; i < philosoperCount; i++ {
		forkId, forkSignal := junction.NewAsyncSignal[types.Unit](j)
		forks = append(forks, Fork{
			Id:   forkId,
			Free: forkSignal,
		})
	}

	// Create Philosophers
	for i := 0; i < philosoperCount; i++ {
		fmt.Println("Setting up philosopher ", i)
		philosopherId, philosopherSignal := junction.NewAsyncSignal[types.Unit](j)
		philosophers = append(philosophers, Philosopher{
			Id:  philosopherId,
			Eat: philosopherSignal,
		})

		func(philosopher int) {
			junction.NewTernaryAsyncJoinPattern[types.Unit, types.Unit, types.Unit](j, forks[philosopher].Id, forks[(philosopher+1)%philosoperCount].Id, philosopherId).
				Action(func(a types.Unit, b types.Unit, c types.Unit) {
					fmt.Println("philosopher", philosopher, "is eating")
					time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
					forks[philosopher].Free(types.Unit{})
					forks[(philosopher+1)%philosoperCount].Free(types.Unit{})

					time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
					philosophers[philosopher].Eat(types.Unit{})
				})
		}(i)
	}

	for _, fork := range forks {
		fork.Free(types.Unit{})
	}

	for i, philosopher := range philosophers {
		fmt.Println(i, " wants to eat")
		go philosopher.Eat(types.Unit{})
	}

	for true {
	}
}
