package helper

import (
	"../types"
)

// CheckForSameJunction checks if all Ports are registered with the same junction
func CheckForSameJunction(ports []types.Port) bool {
	for i := 1; i < len(ports); i++ {
		if ports[i].JunctionChannel != ports[0].JunctionChannel {
			return false
		}
	}

	return true
}
