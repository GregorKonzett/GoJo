package gojo

import (
	"./patterns/binary"
	"./patterns/unary"
	"./types"
	"errors"
	"fmt"
)

type Junction struct {
	sender     chan types.Packet
	receiver   chan interface{}
	JunctionId int
}

func NewJunction() *Junction {
	sender := make(chan types.Packet)
	receiver := make(chan interface{})

	StartController(sender, receiver)

	return &Junction{sender, receiver, 1}
}

func NewAsyncSignal[T any](j *Junction) (types.SignalId, func(T)) {
	channel := getNewChannelId(j)

	return types.SignalId{
			ChannelType: types.AsyncSignal,
			Id:          channel,
			JunctionId:  (*j).JunctionId,
		}, func(data T) {
			fmt.Println("Sending from channel: ", channel)
			(*j).sender <- types.Packet{
				Type: types.MESSAGE,
				Msg:  data,
			}
		}
}

func NewSyncSignal[T any](j *Junction) (types.SignalId, func() (T, error)) {
	channel := getNewChannelId(j)

	return types.SignalId{
			ChannelType: types.SyncSignal,
			Id:          channel,
			JunctionId:  (*j).JunctionId,
		}, func() (T, error) {
			fmt.Println("Sending from channel: ", channel)
			recvChannel := make(chan interface{})

			(*j).sender <- types.Packet{
				Type: types.MESSAGE,
				Ch:   recvChannel,
			}

			receivedData := <-recvChannel

			var returnData T

			switch t := receivedData.(type) {
			case T:
				returnData := t
				return returnData, nil
			default:
				return returnData, errors.New("invalid data type")
			}
		}
}

func NewBiDirSyncSignal[T any, R any](j *Junction) (types.SignalId, func(T) (R, error)) {
	channel := getNewChannelId(j)

	return types.SignalId{
			ChannelType: types.BiDirSyncSignal,
			Id:          channel,
			JunctionId:  (*j).JunctionId,
		}, func(data T) (R, error) {
			fmt.Println("Sending from channel: ", channel)
			recvChannel := make(chan interface{})

			(*j).sender <- types.Packet{
				Type: types.MESSAGE,
				Msg:  data,
				Ch:   recvChannel,
			}

			receivedData := <-recvChannel

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
	(*j).sender <- types.Packet{Type: types.GetNewChannelId}
	channelId := <-(*j).receiver

	switch t := channelId.(type) {
	case int:
		return t
	}

	return 0
}

func NewUnarySendJoinPattern[T any](j *Junction, signal types.SignalId) unary.SendPartialPattern[T] {
	return unary.SendPartialPattern[T]{
		Port:       (*j).sender,
		SignalOne:  signal,
		JunctionId: (*j).JunctionId,
	}
}

func NewBinarySendJoinPattern[T any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId) binary.SendPartialPattern[T, R] {
	return binary.SendPartialPattern[T, R]{
		Port:       (*j).sender,
		JunctionId: (*j).JunctionId,
		SignalOne:  signalOne,
		SignalTwo:  signalTwo,
	}
}

func NewBinaryRecvJoinPattern[T any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId) binary.RecvPartialPattern[T, R] {
	return binary.RecvPartialPattern[T, R]{
		Port:       (*j).sender,
		JunctionId: (*j).JunctionId,
		SignalOne:  signalOne,
		SignalTwo:  signalTwo,
	}
}
