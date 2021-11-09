package helper

import (
	"../types"
)

func CheckForSameJunction(ports []types.SignalId) bool {
	for i := 1; i < len(ports); i++ {
		if ports[i].Junction != ports[0].Junction {
			return false
		}
	}

	return true
}
