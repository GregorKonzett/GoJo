package gojo

import (
	"./types"
	"fmt"
)

func StartController(receiver chan types.Packet, sender chan interface{}) {
	go runThread(receiver, sender)
}

func runThread(receiver chan types.Packet, sender chan interface{}) {
	channelIds := 0

	for true {
		data := <-receiver

		switch data.Type {
		case types.MESSAGE:
			fmt.Println("Incoming message: ", data.Msg)

			if data.Ch != nil {
				data.Ch <- "Message processed"
			}
		case types.AddJoinPattern:
			fmt.Println("Adding new join pattern: ", data.Msg)
			sender <- "Added new join pattern"
		case types.GetNewChannelId:
			fmt.Println("Getting new channel id: ", channelIds)
			sender <- channelIds
			channelIds++
		}
	}
}
