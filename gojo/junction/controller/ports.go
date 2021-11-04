package controller

import "fmt"
import "../../types"

func getNewPortId(patterns *JoinPatterns, msg types.Packet) {
	fmt.Println("Getting new port id: ", (*patterns).portIds)
	msg.Payload.Ch <- (*patterns).portIds
	(*patterns).portIds++
}
