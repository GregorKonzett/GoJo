package types

// Action is used to specify a Packet's type
type Action int

// Unit is introduced for Signals that only need to send void data to a Port
type Unit struct{}

// Defines the different Message Status values
const (
	PENDING  = iota
	CLAIMED  = iota
	CONSUMED = iota
)

// Defines the types of Packets that can be sent to the controller
const (
	MESSAGE        = iota
	AddJoinPattern = iota
	CreateNewPort  = iota
	Shutdown       = iota
)

// Packet is the struct sent to the controller containing different kinds of payloads depending on the type
type Packet struct {
	PortId  int
	Type    Action
	Payload Payload
}

// PortCreation is the controller's response when a new port is created. It contains the channel that's used to send
//messages on, and it's PortId that is registered with join patterns
type PortCreation struct {
	Ch     chan Packet
	PortId int
}

// Payload contains a Msg (which depends on the Payload type), an optional channel that's used to respond
//values to the sender and a Status field (PENDING, CLAIMED, CONSUMED) that's atomically swapped during the pattern matching algorithm
type Payload struct {
	Msg      interface{}
	Ch       chan interface{}
	Status   int32
	Consumed bool
}

//Port combines the PortId and the channel used to communicate with the junction's controller that manages this Port
type Port struct {
	Id              int
	JunctionChannel chan Packet
}

//JoinPatternPacket is sent to the junction's controller when a new join pattern is registered, and it contains the
//ports the pattern is listening on + the function that will be called once the pattern is fired
type JoinPatternPacket struct {
	Ports  []Port
	Action interface{}
}

//UnaryAsync Syntax Sugar to encapsulate the wrapped Action function/**
type UnaryAsync = func(interface{})

//UnarySync Syntax Sugar to encapsulate the wrapped Action function/**
type UnarySync = func(interface{}) interface{}

//BinaryAsync Syntax Sugar to encapsulate the wrapped Action function/**
type BinaryAsync = func(interface{}, interface{})

//BinarySync Syntax Sugar to encapsulate the wrapped Action function/**
type BinarySync = func(interface{}, interface{}) interface{}

//TernaryAsync Syntax Sugar to encapsulate the wrapped Action function/**
type TernaryAsync = func(interface{}, interface{}, interface{})

//TernarySync Syntax Sugar to encapsulate the wrapped Action function/**
type TernarySync = func(interface{}, interface{}, interface{}) interface{}
