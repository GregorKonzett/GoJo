package gojo

type Message[T any] struct {
	Data T
}

type Packet[T any] struct {
	Msg Message[T]
}
