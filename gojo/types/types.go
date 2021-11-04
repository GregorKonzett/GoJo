package types

type Action int
type SignalType int
type Unit struct{}

const (
	AsyncSignal = iota
	SyncSignal  = iota
)

const (
	MESSAGE          = iota
	AddJoinPattern   = iota
	GetNewPortId     = iota
	GetNewJunctionId = iota
)

type Packet struct {
	Port SignalId
	Msg  interface{}
	Type Action
	Ch   chan interface{}
}

type SignalId struct {
	ChannelType SignalType
	Id          int
	JunctionId  int
}

type JoinPatternPacket struct {
	InputPorts  []SignalId
	OutputPorts []SignalId
	DoFunction  interface{}
}

type UnaryAsync func(interface{})
type UnarySync func() interface{}

type BinaryASync func(interface{}, interface{})
type BinarySync func(interface{}) interface{}

type TernaryAsync func(interface{}, interface{}, interface{})
type TernarySync func(interface{}, interface{}) interface{}
