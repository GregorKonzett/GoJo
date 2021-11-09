package helper

import (
	"../types"
)

func CheckForSameJunction(ports []types.Port) bool {
	for i := 1; i < len(ports); i++ {
		if ports[i].JunctionChannel != ports[0].JunctionChannel {
			return false
		}
	}

	return true
}
