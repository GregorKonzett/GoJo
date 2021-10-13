package gojo

type Action int

const (
	ADD_JUNCTION = iota
	MESSAGE      = iota
	ADD_CHANNEL  = iota
)

type Packet struct {
	Msg     interface{}
	Channel chan interface{}
	Type    Action
}

type GenericChan[T any] chan T
