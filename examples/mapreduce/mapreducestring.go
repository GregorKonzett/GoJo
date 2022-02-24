package main

import (
	"../../gojo/junction"
	"../../gojo/types"
	"fmt"
	"strings"
)

type Result struct {
	counts map[string]int
	amount int
}

func createMapReduce1(vals []string) func(types.Unit) (map[string]int, error) {
	j := junction.NewJunction()

	mapPort, mapSignal := junction.NewAsyncPort[string](j)
	reducePort, reduceSignal := junction.NewAsyncPort[map[string]int](j)
	listPort, listSignal := junction.NewAsyncPort[Result](j)
	finalPort, finalSignal := junction.NewAsyncPort[map[string]int](j)
	getPort, getSignal := junction.NewSyncPort[types.Unit, map[string]int](j)

	// map
	junction.NewUnaryAsyncJoinPattern[string](mapPort).Action(func(val string) {
		fmt.Println("Mapping: ", val)
		res := make(map[string]int)

		for _, v := range strings.Split(val, " ") {
			if _, ok := res[v]; !ok {
				res[v] = 1
			} else {
				res[v] = res[v] + 1
			}
		}

		reduceSignal(res)
	})

	// reduce
	junction.NewBinaryAsyncJoinPattern[Result, map[string]int](listPort, reducePort).Action(func(res Result, val map[string]int) {
		fmt.Println("Reducing: ", val)
		for k, v := range val {
			if _, ok := res.counts[k]; !ok {
				res.counts[k] = v
			} else {
				res.counts[k] = res.counts[k] + v
			}
		}

		if res.amount != len(vals)-1 {
			res.amount = res.amount + 1
			listSignal(res)
		} else {
			finalSignal(res.counts)
		}
	})

	// final
	junction.NewBinarySyncJoinPattern[map[string]int, types.Unit, map[string]int](finalPort, getPort).Action(func(res map[string]int, a types.Unit) map[string]int {
		return res
	})

	listSignal(Result{
		counts: make(map[string]int),
		amount: 0,
	})

	for _, val := range vals {
		go mapSignal(val)
	}

	return getSignal
}

func main() {

	str := "hello world\nbye world\nhello map reduce"

	get := createMapReduce1(strings.Split(str, "\n"))
	items, _ := get(types.Unit{})

	fmt.Println("Result: ", items)
}
