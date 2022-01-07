package main

import (
	"../../gojo/junction"
	"../../gojo/types"
	"bytes"
	"constraints"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"strings"
)

type MapReduce[T any, U constraints.Ordered, Z any] struct {
	mapperSignals   []func(T)
	combinerSignals []func(map[U]Z)
}

type ResultContainer[U constraints.Ordered, Z any] struct {
	counts map[U]Z
	amount int
}

func hash[T any](s T) uint32 {
	h := fnv.New32a()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	enc.Encode(s)
	h.Write(buf.Bytes())

	return h.Sum32()
}

// Creates a map reduce instance with n mapper and m combiner
func createMapReduce2[T any, U constraints.Ordered, Z any](
	n int,
	m int,
	mapFunc func(T) map[U]Z,
	combinerFunc func(map[U]Z, map[U]Z),
	reduceFunc func(ResultContainer[U, Z], map[U]Z),
	data []T) func(types.Unit) (map[U]Z, error) {

	mr := MapReduce[T, U, Z]{
		mapperSignals:   make([]func(T), 0),
		combinerSignals: make([]func(map[U]Z), 0),
	}

	reducerJunction := junction.NewJunction()

	reducePort, reduceSignal := junction.NewAsyncSignal[map[U]Z](reducerJunction)
	listPort, listSignal := junction.NewAsyncSignal[ResultContainer[U, Z]](reducerJunction)
	finalPort, finalSignal := junction.NewAsyncSignal[map[U]Z](reducerJunction)
	getPort, getSignal := junction.NewSyncSignal[types.Unit, map[U]Z](reducerJunction)

	// Create 1 final reducer merging all combiner outputs
	junction.NewBinaryAsyncJoinPattern[ResultContainer[U, Z], map[U]Z](listPort, reducePort).Action(func(res ResultContainer[U, Z], val map[U]Z) {
		reduceFunc(res, val)

		if res.amount != 1 {
			res.amount = res.amount - 1
			listSignal(res)
		} else {
			finalSignal(res.counts)
		}
	})

	// final join pattern returning result to caller
	junction.NewBinarySyncJoinPattern[map[U]Z, types.Unit, map[U]Z](finalPort, getPort).Action(func(res map[U]Z, a types.Unit) map[U]Z {
		return res
	})

	// Creating m combiner
	for i := 0; i < m; i++ {
		combinerJunction := junction.NewJunction()
		combinerPort, combinerSignal := junction.NewAsyncSignal[map[U]Z](combinerJunction)
		combinerResPort, combinerResSignal := junction.NewAsyncSignal[ResultContainer[U, Z]](combinerJunction)
		mr.combinerSignals = append(mr.combinerSignals, combinerSignal)

		junction.NewBinaryAsyncJoinPattern[ResultContainer[U, Z], map[U]Z](combinerResPort, combinerPort).Action(func(res ResultContainer[U, Z], val map[U]Z) {
			combinerFunc(res.counts, val)

			if res.amount == 1 {
				// Send to Reducer
				reduceSignal(res.counts)
			} else {
				combinerResSignal(ResultContainer[U, Z]{
					counts: res.counts,
					amount: res.amount - 1,
				})
			}
		})

		combinerResSignal(ResultContainer[U, Z]{
			counts: make(map[U]Z),
			amount: len(data),
		})
	}

	// Creating n Mapper
	for i := 0; i < n; i++ {
		mapperPort, mapperSignal := junction.NewAsyncSignal[T](junction.NewJunction())
		mr.mapperSignals = append(mr.mapperSignals, mapperSignal)

		junction.NewUnaryAsyncJoinPattern[T](mapperPort).Action(func(val T) {
			result := mapFunc(val)

			combinerCalls := make(map[uint32]map[U]Z)
			var usedCombiners []uint32

			for k, v := range result {
				h := hash(k) % uint32(len(mr.combinerSignals))
				if _, ok := combinerCalls[h]; !ok {
					combinerCalls[h] = make(map[U]Z)
				}

				combinerCalls[h][k] = v
				usedCombiners = append(usedCombiners, h)
			}

			// Send each mapped result to a combiner
			for k, v := range combinerCalls {
				mr.combinerSignals[k](v)
			}

			for i, _ := range mr.combinerSignals {
				if !contains(usedCombiners, uint32(i)) {
					mr.combinerSignals[i](make(map[U]Z))
				}
			}
		})
	}

	listSignal(ResultContainer[U, Z]{
		counts: make(map[U]Z),
		amount: len(mr.combinerSignals),
	})

	process[T, U, Z](mr, data)

	return getSignal
}

func process[T any, U constraints.Ordered, Z any](mr MapReduce[T, U, Z], data []T) {
	for i := 0; i < len(data); i++ {
		mr.mapperSignals[i%len(mr.mapperSignals)](data[i])
	}
}

func contains(combiners []uint32, combiner uint32) bool {
	for _, elem := range combiners {
		if elem == combiner {
			return true
		}
	}
	return false
}

func main() {
	data := "hello world\nhello reduce\nmap reduce\nnothing here here"
	getResult := createMapReduce2[string, string, int](2, 3, func(val string) map[string]int {
		res := make(map[string]int)

		for _, v := range strings.Split(val, " ") {
			if _, ok := res[v]; !ok {
				res[v] = 1
			} else {
				res[v] = res[v] + 1
			}
		}

		return res
	}, func(res map[string]int, vals map[string]int) {
		for k, v := range vals {
			if _, ok := res[k]; !ok {
				res[k] = v
			} else {
				res[k] = res[k] + v
			}
		}
	}, func(res ResultContainer[string, int], val map[string]int) {
		for k, v := range val {
			if _, ok := res.counts[k]; !ok {
				res.counts[k] = v
			} else {
				res.counts[k] = res.counts[k] + v
			}
		}
	}, strings.Split(data, "\n"))

	val, _ := getResult(types.Unit{})

	fmt.Println("Done: ", val)
}
