package controller

import (
	"../../types"
	"fmt"
)

func registerNewJoinPattern(patterns *JoinPatterns, joinPattern types.JoinPatternPacket) {
	fmt.Println("Adding new join pattern: ")
	(*patterns).joinPatterns[(*patterns).joinPatternId] = joinPattern

	fillPortsToJoinPatterns(patterns, joinPattern.Signals)

	(*patterns).joinPatternId++
}

func fillPortsToJoinPatterns(patterns *JoinPatterns, signals []types.SignalId) {
	for _, port := range signals {
		if (*patterns).portsToJoinPatterns[port.Id] == nil {
			(*patterns).portsToJoinPatterns[port.Id] = []int{}
		}

		(*patterns).portsToJoinPatterns[port.Id] = append((*patterns).portsToJoinPatterns[port.Id], (*patterns).joinPatternId)
	}
}
