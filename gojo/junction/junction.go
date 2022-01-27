package junction

import (
	"../patterns/binary"
	"../patterns/ternary"
	"../patterns/unary"
	"../types"
	"./controller"
	"errors"
)

type Junction struct {
	port chan types.Packet
}

func NewJunction() *Junction {
	sender := make(chan types.Packet)

	controller.StartController(sender)

	return &Junction{sender}
}

func NewAsyncSignal[T any](j *Junction) (types.Port, func(T)) {
	portNr, signalChannel := createNewPort(j)

	signalId := types.Port{
		Id:              portNr,
		JunctionChannel: (*j).port,
	}

	return signalId, func(data T) {
		signalChannel <- &types.Payload{
			Msg: data,
		}
	}
}

func NewSyncSignal[T any, R any](j *Junction) (types.Port, func(T) (R, error)) {
	portNr, signalChannel := createNewPort(j)

	signalId := types.Port{
		Id:              portNr,
		JunctionChannel: (*j).port,
	}

	return signalId, func(data T) (R, error) {
		recvChannel := make(chan interface{})

		signalChannel <- &types.Payload{
			Msg: data,
			Ch:  recvChannel,
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

func Shutdown(j *Junction) {
	(*j).port <- types.Packet{Type: types.Shutdown}
}

func createNewPort(j *Junction) (int, chan *types.Payload) {
	receiver := make(chan interface{})
	(*j).port <- types.Packet{Type: types.CreateNewPort, Payload: types.Payload{Ch: receiver}}
	signalChannel := <-receiver

	switch t := signalChannel.(type) {
	case types.PortCreation:
		return t.SignalId, t.Ch
	}

	return 0, nil
}

func NewUnaryAsyncJoinPattern[T any](signal types.Port) unary.AsyncPartialPattern[T] {
	return unary.AsyncPartialPattern[T]{
		Signals: []types.Port{signal},
	}
}

func NewUnarySyncJoinPattern[T any, R any](signal types.Port) unary.SyncPartialPattern[T, R] {
	return unary.SyncPartialPattern[T, R]{
		Signals: []types.Port{signal},
	}
}

func NewBinaryAsyncJoinPattern[T any, R any](signalOne types.Port, signalTwo types.Port) binary.AsyncPartialPattern[T, R] {
	return binary.AsyncPartialPattern[T, R]{
		Signals: []types.Port{signalOne, signalTwo},
	}
}

func NewBinarySyncJoinPattern[T any, S any, R any](signalOne types.Port, signalTwo types.Port) binary.SyncPartialPattern[T, S, R] {
	return binary.SyncPartialPattern[T, S, R]{
		Signals: []types.Port{signalOne, signalTwo},
	}
}

func NewTernaryAsyncJoinPattern[T any, S any, R any](signalOne types.Port, signalTwo types.Port, signalThree types.Port) ternary.AsyncPartialPattern[T, S, R] {
	return ternary.AsyncPartialPattern[T, S, R]{
		Signals: []types.Port{signalOne, signalTwo, signalThree},
	}
}

func NewTernarySyncJoinPattern[T any, S any, R any, U any](signalOne types.Port, signalTwo types.Port, signalThree types.Port) ternary.SyncPartialPattern[T, S, R, U] {
	return ternary.SyncPartialPattern[T, S, R, U]{
		Signals: []types.Port{signalOne, signalTwo, signalThree},
	}
}
