package controller

import (
	"../../types"
	"fmt"
)

func registerNewJoinPattern(joinPatterns *JoinPatterns, joinPattern types.JoinPatternPacket) {
	fmt.Println("Adding new join pattern: ")
	(*joinPatterns).registeredJoinPatterns[(*joinPatterns).joinPatternId] = joinPattern

	for _, port := range joinPattern.InputPorts {
		(*joinPatterns).firedPorts[port.Id] = 0
		(*joinPatterns).portIdToJoinPatternId[port.Id] = (*joinPatterns).joinPatternId
	}

	for _, port := range joinPattern.OutputPorts {
		(*joinPatterns).firedPorts[port.Id] = 0
		(*joinPatterns).portIdToJoinPatternId[port.Id] = (*joinPatterns).joinPatternId
	}

	(*joinPatterns).joinPatternId++
}
