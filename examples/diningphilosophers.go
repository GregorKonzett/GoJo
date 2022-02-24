package main

import (
	"../gojo/junction"
	"../gojo/types"
	"fmt"
	"math/rand"
	"time"
)

type Fork struct {
	Id   types.Port
	Free func(types.Unit)
}

type Philosopher struct {
	Id  types.Port
	Eat func(types.Unit)
}

func main() {
	philosopherCount := 5
	j := junction.NewJunction()

	var forks []Fork
	var philosophers []Philosopher

	// Create Forks
	for i := 0; i < philosopherCount; i++ {
		forkId, forkSignal := junction.NewAsyncPort[types.Unit](j)
		forks = append(forks, Fork{
			Id:   forkId,
			Free: forkSignal,
		})
	}

	// Create Philosophers
	for i := 0; i < philosopherCount; i++ {
		fmt.Println("Setting up philosopher ", i)
		philosopherId, eat := junction.NewAsyncPort[types.Unit](j)
		sleepId, sleepSignal := junction.NewAsyncPort[types.Unit](j)
		philosophers = append(philosophers, Philosopher{
			Id:  philosopherId,
			Eat: eat,
		})

		func(philosopher int, sleepId types.Port, sleep func(types.Unit)) {
			junction.NewTernaryAsyncJoinPattern[types.Unit, types.Unit, types.Unit](forks[philosopher].Id, forks[(philosopher+1)%philosopherCount].Id, philosopherId).
				Action(func(a types.Unit, b types.Unit, c types.Unit) {
					fmt.Println("philosopher", philosopher, "is eating")

					forks[philosopher].Free(types.Unit{})
					forks[(philosopher+1)%philosopherCount].Free(types.Unit{})
					sleep(types.Unit{})
				})

			junction.NewUnaryAsyncJoinPattern[types.Unit](sleepId).Action(func(a types.Unit) {
				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				fmt.Println(philosopher, " wants to eat")
				philosophers[philosopher].Eat(types.Unit{})
			})
		}(i, sleepId, sleepSignal)
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
