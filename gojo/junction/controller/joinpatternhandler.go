package controller

import (
	"../../types"
)

func registerNewJoinPattern(patterns *JoinPatterns, joinPattern types.JoinPatternPacket) {
	(*patterns).joinPatterns[(*patterns).joinPatternId] = joinPattern

	fillPortsToJoinPatterns(patterns, joinPattern.Signals)

	(*patterns).joinPatternId++
}

func fillPortsToJoinPatterns(patterns *JoinPatterns, signals []types.Port) {
	for _, port := range signals {
		if (*patterns).portsToJoinPatterns[port.Id] == nil {
			(*patterns).portsToJoinPatterns[port.Id] = []int{}
		}

		(*patterns).portsToJoinPatterns[port.Id] = append((*patterns).portsToJoinPatterns[port.Id], (*patterns).joinPatternId)
	}
}
