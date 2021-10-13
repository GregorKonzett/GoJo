package gojo

import (
	"fmt"
)

func StartController(receiver chan Packet, sender chan interface{}) {
	go runThread(receiver, sender)
}

func runThread(receiver chan Packet, sender chan interface{}) {
	channelIds := 0

	for true {
		data := <-receiver

		switch data.Type {
		case MESSAGE:
			fmt.Println("Incoming message: ", data.Msg)
			sender <- "Message processed"
		case AddJoinPattern:
			fmt.Println("Adding new join pattern: ", data.Msg)
			sender <- "Added new join pattern"
		case GetNewChannelId:
			fmt.Println("Getting new channel id: ", channelIds)
			sender <- channelIds
			channelIds++
		}
	}
}
