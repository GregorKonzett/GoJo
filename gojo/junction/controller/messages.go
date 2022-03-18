package controller

import (
	"../../types"
	"sync/atomic"
)

// tryClaimMessages is a greedy algorithm trying to claim one message from each of the join pattern's ports to fire the
// join pattern. Since Payloads are shared across all join patterns the Payloads are consumed atomically using
// compare and swap
func tryClaimMessages(params map[int][]*types.Packet, portOrders []int) ([]interface{}, []chan interface{}, bool) {
	retry := true

	messages := make([]interface{}, len(portOrders))
	var syncPorts []chan interface{}

	for retry {
		var chosenParams []*types.Payload

		for _, portId := range portOrders {
			foundPending := findPending(params[portId])

			if foundPending == nil {
				for _, param := range chosenParams {
					(*param).Consumed = false
				}

				return nil, nil, false
			}

			chosenParams = append(chosenParams, foundPending)
		}

		alreadyConsumedParams := make([]bool, len(chosenParams))
		releaseParams := false

		for i, chosenParam := range chosenParams {
			if !atomic.CompareAndSwapInt32(&chosenParam.Status, types.PENDING, types.CLAIMED) {
				alreadyConsumedParams[i] = true
				releaseParams = true
			}
		}

		if releaseParams {
			for i, chosenParam := range chosenParams {
				if !alreadyConsumedParams[i] {
					chosenParam.Status = types.PENDING
					chosenParam.Consumed = false
				}
			}
			continue
		}

		for i, chosenParam := range chosenParams {
			chosenParam.Status = types.CONSUMED
			messages[i] = (*chosenParam).Msg

			if (*chosenParam).Ch != nil {
				syncPorts = append(syncPorts, (*chosenParam).Ch)
			}
		}

		retry = false
	}

	return messages, syncPorts, true
}

// findPending loops through every available WrappedPayload to find the first PENDING Payload that hasn't been consumed
// yet during this iteration of the pattern matching algorithm. This ensures that a message won't be consumed twice if
// a join pattern listens on the same port multiple times (non linear join pattern)
func findPending(params []*types.Packet) *types.Payload {
	for _, param := range params {
		if (*param).Payload.Consumed {
			continue
		}

		if (*param).Payload.Status == types.PENDING {
			(*param).Payload.Consumed = true
			return &(*param).Payload
		}
	}

	return nil
}
