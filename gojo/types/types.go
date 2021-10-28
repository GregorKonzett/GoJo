package types

type Action int
type SignalType int

const (
	AsyncSignal     = iota
	SyncSignal      = iota
	BiDirSyncSignal = iota
)

const (
	MESSAGE         = iota
	AddJoinPattern  = iota
	GetNewChannelId = iota
)

type Packet struct {
	Msg  interface{}
	Type Action
	Ch   chan interface{}
}

type SignalId struct {
	ChannelType SignalType
	Id          int
	JunctionId  int
}
