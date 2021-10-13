package gojo

type Action int

const (
	MESSAGE         = iota
	AddJoinPattern  = iota
	GetNewChannelId = iota
)

type Packet struct {
	Msg  interface{}
	Type Action
}
