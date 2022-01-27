package types

type Action int
type SignalType int
type Unit struct{}

const (
	PENDING  = iota
	CLAIMED  = iota
	CONSUMED = iota
)

const (
	MESSAGE        = iota
	AddJoinPattern = iota
	CreateNewPort  = iota
	Shutdown       = iota
)

type Packet struct {
	SignalId Port
	Type     Action
	Payload  Payload
}

type PortCreation struct {
	Ch       chan *Payload
	SignalId int
}

type Payload struct {
	Msg    interface{}
	Ch     chan interface{}
	status int
}

type WrappedPayload struct {
	Payload *Payload
	PortId  int
}

type Port struct {
	Id              int
	JunctionChannel chan Packet
}

type JoinPatternPacket struct {
	Ports  []Port
	Action interface{}
}

type WrappedJoinPattern struct {
	Pattern JoinPatternPacket
	Bitmask int
}

type UnaryAsync = func(interface{})
type UnarySync = func(interface{}) interface{}

type BinaryAsync = func(interface{}, interface{})
type BinarySync = func(interface{}, interface{}) interface{}

type TernaryAsync = func(interface{}, interface{}, interface{})
type TernarySync = func(interface{}, interface{}, interface{}) interface{}
