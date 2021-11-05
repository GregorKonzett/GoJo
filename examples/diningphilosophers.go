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
	Eat func(int)
}

func main() {
	j := junction.NewJunction()

	var forks []Fork
	var philosophers []Philosopher

	// Create Forks
	for i := 0; i < 5; i++ {
		forkId, forkSignal := junction.NewAsyncSignal[types.Unit](j)
		forks = append(forks, Fork{
			Id:   forkId,
			Free: forkSignal,
		})
	}

	// Create Philosophers
	for i := 0; i < 5; i++ {
		fmt.Println("Setting up philosopher ", i)
		philosopherId, philosopherSignal := junction.NewAsyncSignal[int](j)
		philosophers = append(philosophers, Philosopher{
			Id:  philosopherId,
			Eat: philosopherSignal,
		})

		junction.NewTernaryAsyncJoinPattern[types.Unit, types.Unit, int](j, forks[i].Id, forks[(i+1)%5].Id, philosopherId).
			ThenDo(func(a types.Unit, b types.Unit, c int) {
				fmt.Println("philosopher", c, "is eating")
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
				forks[c].Free(types.Unit{})
				forks[(c+1)%5].Free(types.Unit{})

				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				philosophers[c].Eat(c)
			})
	}

	for _, fork := range forks {
		fork.Free(types.Unit{})
	}

	for i, philosopher := range philosophers {
		fmt.Println(i, " wants to eat")
		go philosopher.Eat(i)
	}

	for true {
	}
}
