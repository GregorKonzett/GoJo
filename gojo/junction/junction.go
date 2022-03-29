package junction

import (
	"errors"
	"github.com/junctional/GoJo/gojo/junction/controller"
	"github.com/junctional/GoJo/gojo/patterns/binary"
	"github.com/junctional/GoJo/gojo/patterns/ternary"
	"github.com/junctional/GoJo/gojo/patterns/unary"
	"github.com/junctional/GoJo/gojo/types"
)

// Junction is the entrypoint to register new Signals and Join Patterns. The only necessary information is the port to
// the controller goroutine handling the registration processes
type Junction struct {
	port chan types.Packet
}

// NewJunction creates a new Junction, starts the controller goroutine in the background and returns a reference to
// this junction
func NewJunction() *Junction {
	sender := make(chan types.Packet)

	controller.StartController(sender)

	return &Junction{sender}
}

// NewAsyncPort Creates a new Port,Signal pair by registering a new Port on the controller goroutine. This Signal will
// not receive a return value.
func NewAsyncPort[T any](j *Junction) (types.Port, func(T)) {
	portNr, signalChannel := createNewPort(j)

	portId := types.Port{
		Id:              portNr,
		JunctionChannel: (*j).port,
	}

	return portId, func(data T) {
		signalChannel <- types.Packet{
			PortId: portNr,
			Type:   types.MESSAGE,
			Payload: types.Payload{
				Msg:    data,
				Status: types.PENDING,
			},
		}
	}
}

// NewSyncPort Creates a new Port,Signal pair by registering a new Port on the controller goroutine. This Signal will
// receive a return value and will block until it receives a value.
func NewSyncPort[T any, R any](j *Junction) (types.Port, func(T) (R, error)) {
	portNr, signalChannel := createNewPort(j)

	portId := types.Port{
		Id:              portNr,
		JunctionChannel: (*j).port,
	}

	return portId, func(data T) (R, error) {
		recvChannel := make(chan interface{})

		signalChannel <- types.Packet{
			PortId: portNr,
			Type:   types.MESSAGE,
			Payload: types.Payload{
				Msg:    data,
				Ch:     recvChannel,
				Status: types.PENDING,
			},
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

// Shutdown stops the controller goroutine and the junction
func Shutdown(j *Junction) {
	(*j).port <- types.Packet{Type: types.Shutdown}
}

// createNewPort sends a Packet to the controller goroutine and returning it's Port ID + the channel to Signal
func createNewPort(j *Junction) (int, chan types.Packet) {
	receiver := make(chan interface{})
	(*j).port <- types.Packet{Type: types.CreateNewPort, Payload: types.Payload{Ch: receiver}}
	signalChannel := <-receiver

	switch t := signalChannel.(type) {
	case types.PortCreation:
		return t.PortId, t.Ch
	}

	return 0, nil
}

// NewUnaryAsyncJoinPattern takes a Port and the Action data type to create a new Join Pattern, which still
// requires an assigned Action to be registered at the controller.
func NewUnaryAsyncJoinPattern[T any](signal types.Port) unary.AsyncPattern[T] {
	return unary.AsyncPattern[T]{
		Signals: []types.Port{signal},
	}
}

// NewUnarySyncJoinPattern takes a list of Ports and the Action data types to create a new Join Pattern, which still
// requires an assigned Action to be registered at the controller.
func NewUnarySyncJoinPattern[T any, R any](signal types.Port) unary.SyncPattern[T, R] {
	return unary.SyncPattern[T, R]{
		Signals: []types.Port{signal},
	}
}

// NewBinaryAsyncJoinPattern takes a list of Ports and the Action data types to create a new Join Pattern, which still
// requires an assigned Action to be registered at the controller.
func NewBinaryAsyncJoinPattern[T any, R any](signalOne types.Port, signalTwo types.Port) binary.AsyncPattern[T, R] {
	return binary.AsyncPattern[T, R]{
		Signals: []types.Port{signalOne, signalTwo},
	}
}

// NewBinarySyncJoinPattern takes a list of Ports and the Action data types to create a new Join Pattern, which still
// requires an assigned Action to be registered at the controller.
func NewBinarySyncJoinPattern[T any, S any, R any](signalOne types.Port, signalTwo types.Port) binary.SyncPattern[T, S, R] {
	return binary.SyncPattern[T, S, R]{
		Signals: []types.Port{signalOne, signalTwo},
	}
}

// NewTernaryAsyncJoinPattern takes a list of Ports and the Action data types to create a new Join Pattern, which still
// requires an assigned Action to be registered at the controller.
func NewTernaryAsyncJoinPattern[T any, S any, R any](signalOne types.Port, signalTwo types.Port, signalThree types.Port) ternary.AsyncPattern[T, S, R] {
	return ternary.AsyncPattern[T, S, R]{
		Signals: []types.Port{signalOne, signalTwo, signalThree},
	}
}

// NewTernarySyncJoinPattern takes a list of Ports and the Action data types to create a new Join Pattern, which still
// requires an assigned Action to be registered at the controller.
func NewTernarySyncJoinPattern[T any, S any, R any, U any](signalOne types.Port, signalTwo types.Port, signalThree types.Port) ternary.SyncPattern[T, S, R, U] {
	return ternary.SyncPattern[T, S, R, U]{
		Signals: []types.Port{signalOne, signalTwo, signalThree},
	}
}
