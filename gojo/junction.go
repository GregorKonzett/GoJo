package gojo

import (
	"errors"
	"fmt"
)

type Junction struct {
	sender   chan interface{}
	receiver chan interface{}
	channels int
}

func NewJunction() *Junction {
	sender := make(chan interface{})
	receiver := make(chan interface{})

	StartController(sender, receiver)

	return &Junction{sender, receiver, 0}
}

func NewAsyncSignal[T any](j *Junction) (int, func(T)) {
	channel := (*j).channels
	(*j).channels++

	return channel, func(data T) {
		fmt.Println("Sending from channel: ", channel)
		(*j).sender <- data
	}
}

func NewSyncSignal[T any, R any](j *Junction) (int, func(T) (R, error)) {
	channel := (*j).channels
	(*j).channels++

	return channel, func(data T) (R, error) {
		fmt.Println("Sending from channel: ", channel)
		(*j).sender <- data

		receivedData := <-(*j).receiver

		var returnData R

		switch t := receivedData.(type) {
		case R:
			returnData := t
			return returnData, nil
		default:
			return returnData, errors.New("invalid data type")
		}
	}
}
