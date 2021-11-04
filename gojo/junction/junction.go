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
	sender     chan types.Packet
	JunctionId int
}

func NewJunction() *Junction {
	sender := make(chan types.Packet)

	controller.StartController(sender)

	return &Junction{sender, 1}
}

func NewAsyncSignal[T any](j *Junction) (types.SignalId, func(T)) {
	portId := types.SignalId{
		ChannelType: types.AsyncSignal,
		Id:          getNewPortId(j),
		JunctionId:  (*j).JunctionId,
	}

	return portId, func(data T) {
		(*j).sender <- types.Packet{
			Port: portId,
			Type: types.MESSAGE,
			Msg:  data,
		}
	}
}

func NewSyncSignal[T any, R any](j *Junction) (types.SignalId, func(T) (R, error)) {
	portId := types.SignalId{
		ChannelType: types.SyncSignal,
		Id:          getNewPortId(j),
		JunctionId:  (*j).JunctionId,
	}

	return portId, func(data T) (R, error) {
		recvChannel := make(chan interface{})

		(*j).sender <- types.Packet{
			Port: portId,
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

func getNewPortId(j *Junction) int {
	receiver := make(chan interface{})
	(*j).sender <- types.Packet{Type: types.GetNewPortId, Ch: receiver}
	portId := <-receiver

	switch t := portId.(type) {
	case int:
		return t
	}

	return 0
}

func NewUnaryAsyncJoinPattern[T any](j *Junction, signal types.SignalId) unary.AsyncPartialPattern[T] {
	return unary.AsyncPartialPattern[T]{
		Port:         (*j).sender,
		InputSignals: []types.SignalId{signal},
		JunctionId:   (*j).JunctionId,
	}
}

func NewUnarySyncJoinPattern[T any](j *Junction, signal types.SignalId) unary.SyncPartialPattern[T] {
	return unary.SyncPartialPattern[T]{
		Port:          (*j).sender,
		OutputSignals: []types.SignalId{signal},
		JunctionId:    (*j).JunctionId,
	}
}

func NewBinaryAsyncJoinPattern[T any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId) binary.AsyncPartialPattern[T, R] {
	return binary.AsyncPartialPattern[T, R]{
		Port:         (*j).sender,
		JunctionId:   (*j).JunctionId,
		InputSignals: []types.SignalId{signalOne, signalTwo},
	}
}

func NewBinarySyncJoinPattern[T any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId) binary.SyncPartialPattern[T, R] {
	return binary.SyncPartialPattern[T, R]{
		Port:          (*j).sender,
		JunctionId:    (*j).JunctionId,
		InputSignals:  []types.SignalId{signalOne},
		OutputSignals: []types.SignalId{signalTwo},
	}
}

func NewTernaryAsyncJoinPattern[T any, S any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId, signalThree types.SignalId) ternary.AsyncPartialPattern[T, S, R] {
	return ternary.AsyncPartialPattern[T, S, R]{
		Port:         (*j).sender,
		JunctionId:   (*j).JunctionId,
		InputSignals: []types.SignalId{signalOne, signalTwo, signalThree},
	}
}

func NewTernarySyncJoinPattern[T any, S any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId, signalThree types.SignalId) ternary.SyncPartialPattern[T, S, R] {
	return ternary.SyncPartialPattern[T, S, R]{
		Port:          (*j).sender,
		JunctionId:    (*j).JunctionId,
		InputSignals:  []types.SignalId{signalOne, signalTwo},
		OutputSignals: []types.SignalId{signalThree},
	}
}
