package gojo

import (
	"errors"
	"fmt"
)

type Junction struct {
	sender   chan Packet
	receiver chan interface{}
	channels int
}

func NewJunction() *Junction {
	sender := make(chan Packet)
	receiver := make(chan interface{})

	StartController(sender, receiver)

	return &Junction{sender, receiver, 0}
}

func NewAsyncSignal[T any](j *Junction) (ChannelId, func(T)) {
	channel := getNewChannelId(j)

	return ChannelId{
			ASYNC_SIGNAL,
			channel,
		}, func(data T) {
			fmt.Println("Sending from channel: ", channel)
			(*j).sender <- Packet{
				Type: MESSAGE,
				Msg:  data,
			}
		}
}

func NewSyncSignal[T any, R any](j *Junction) (ChannelId, func(T) (R, error)) {
	channel := getNewChannelId(j)

	return ChannelId{
			SYNC_SIGNAL,
			channel,
		}, func(data T) (R, error) {
			fmt.Println("Sending from channel: ", channel)
			(*j).sender <- Packet{
				Type: MESSAGE,
				Msg:  data,
			}

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

func getNewChannelId(j *Junction) int {
	(*j).sender <- Packet{Type: GetNewChannelId}
	channelId := <-(*j).receiver

	switch t := channelId.(type) {
	case int:
		return t
	}

	return 0
}
