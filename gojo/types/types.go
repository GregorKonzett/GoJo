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
	SignalId SignalId
	Type     Action
	Payload  Payload
}

type Payload struct {
	Msg interface{}
	Ch  chan interface{}
}

type SignalId struct {
	ChannelType SignalType
	Id          int
	JunctionId  int
}

type JoinPatternPacket struct {
	Signals []SignalId
	Action  interface{}
}

type UnaryAsync = func(interface{})
type UnarySync = func(interface{}) interface{}

type BinaryAsync = func(interface{}, interface{})
type BinarySync = func(interface{}, interface{}) interface{}

type TernaryAsync = func(interface{}, interface{}, interface{})
type TernarySync = func(interface{}, interface{}, interface{}) interface{}
