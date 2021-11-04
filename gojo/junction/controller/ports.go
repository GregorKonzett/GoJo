package controller

import "fmt"
import "../../types"

func getNewPortId(joinPatterns *JoinPatterns, msg types.Packet) {
	fmt.Println("Getting new port id: ", (*joinPatterns).portIds)
	msg.Ch <- (*joinPatterns).portIds
	(*joinPatterns).portIds++
}
