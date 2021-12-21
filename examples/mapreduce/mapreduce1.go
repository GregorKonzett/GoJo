package main

import (
	"../../gojo/junction"
	"fmt"
	"hash/fnv"
)

type MapReduce[T any, U any] struct {
	mapperSignals   []func(T)
	combinerSignals []func(U)
}

func hash[T any](s T) uint32 {
	h := fnv.New32a()
	tmp := interface{}(s)

	switch tmp.(type) {
	case string:
		h.Write([]byte(tmp.(string)))
	}

	return h.Sum32()
}

// Creates a map reduce instance with n mapper and m combiner
func createMapReduce2[T any, U any, Z any](n int, m int, mapFunc func(T) U, combinerFunc func(U) Z) MapReduce[T, U] {
	mr := MapReduce[T, U]{
		mapperSignals:   make([]func(T), 0),
		combinerSignals: make([]func(U), 0),
	}

	// Creating m combiner
	for i := 0; i < m; i++ {
		combinerPort, combinerSignal := junction.NewAsyncSignal[U](junction.NewJunction())
		mr.combinerSignals = append(mr.combinerSignals, combinerSignal)

		junction.NewUnaryAsyncJoinPattern[U](combinerPort).Action(func(val U) {
			result := combinerFunc(val)
		})
	}

	// Creating n Mapper
	for i := 0; i < n; i++ {
		mapperPort, mapperSignal := junction.NewAsyncSignal[T](junction.NewJunction())
		mr.mapperSignals = append(mr.mapperSignals, mapperSignal)

		junction.NewUnaryAsyncJoinPattern[T](mapperPort).Action(func(val T) {
			result := mapFunc(val)

			// Hash result to one of the combiner
			fmt.Println("Prev: ", val, result)
			mr.combinerSignals[hash(result)%uint32(len(mr.combinerSignals))](result)
		})
	}

	return mr
}

func (mr MapReduce[T, U]) process(data []T) {
	for i := 0; i < len(data); i++ {
		mr.mapperSignals[i%len(mr.mapperSignals)](data[i])
	}
}

func main() {
	data := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	mr := createMapReduce2[string, []string, string](2, 3, func(val string) []string {
		return []string{val, val}
	}, func(vals []string) string {
		fmt.Println("Combining ", vals)
		return vals[0]
	})

	mr.process(data)

	for true {
	}
}
