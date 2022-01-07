package controller

import (
	"../../types"
)

func registerNewJoinPattern(patterns *JoinPatterns, joinPattern types.JoinPatternPacket) {
	(*patterns).joinPatterns[(*patterns).joinPatternId] = types.WrappedJoinPattern{
		Pattern: joinPattern,
		Bitmask: getBitmask(joinPattern.Signals),
	}

	fillPortsToJoinPatterns(patterns, joinPattern.Signals)

	(*patterns).joinPatternId++
}

func getBitmask(signals []types.Port) int {
	bitmask := 0

	for _, port := range signals {
		bitmask |= 1 << port.Id
	}

	return bitmask
}

func fillPortsToJoinPatterns(patterns *JoinPatterns, signals []types.Port) {
	for _, port := range signals {
		if (*patterns).portsToJoinPatterns[port.Id] == nil {
			(*patterns).portsToJoinPatterns[port.Id] = []int{}
		}

		(*patterns).portsToJoinPatterns[port.Id] = append((*patterns).portsToJoinPatterns[port.Id], (*patterns).joinPatternId)
	}
}
