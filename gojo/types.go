package gojo

type Action int
type ChannelType int

const (
	AsyncSignal = iota
	SyncSignal  = iota
)

const (
	MESSAGE         = iota
	AddJoinPattern  = iota
	GetNewChannelId = iota
)

type Packet struct {
	Msg  interface{}
	Type Action
}

type ChannelId struct {
	channelType ChannelType
	id          int
}
