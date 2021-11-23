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
	Shutdown         = iota
)

type Packet struct {
	SignalId Port
	Type     Action
	Payload  Payload
}

type Payload struct {
	Msg interface{}
	Ch  chan interface{}
}

type Port struct {
	ChannelType     SignalType
	Id              int
	JunctionChannel chan Packet
}

type JoinPatternPacket struct {
	Signals []Port
	Action  interface{}
}

type UnaryAsync = func(interface{})
type UnarySync = func(interface{}) interface{}

type BinaryAsync = func(interface{}, interface{})
type BinarySync = func(interface{}, interface{}) interface{}

type TernaryAsync = func(interface{}, interface{}, interface{})
type TernarySync = func(interface{}, interface{}, interface{}) interface{}
