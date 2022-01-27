package controller

import (
	"../../types"
	"sync/atomic"
)

func tryClaimMessages(params map[int][]*types.WrappedPayload, portOrders []types.Port) ([]interface{}, []chan interface{}, bool) {
	retry := true

	var chosenParams []*types.WrappedPayload
	var messages []interface{}
	var syncPorts []chan interface{}

	for retry {
		for _, portId := range portOrders {
			foundPending := findPending(params[portId.Id])

			if foundPending == nil {
				for _, param := range chosenParams {
					(*param).Consumed = false
				}

				return nil, nil, false
			}

			chosenParams = append(chosenParams, foundPending)
		}

		paramWasAlreadyConsumed := false

		for _, chosenParam := range chosenParams {
			if !atomic.CompareAndSwapInt32(&(*(*chosenParam).Payload).Status, types.PENDING, types.CLAIMED) {
				paramWasAlreadyConsumed = true
				break
			}
		}

		if paramWasAlreadyConsumed {
			for _, chosenParam := range chosenParams {
				(*(*chosenParam).Payload).Status = types.PENDING
				(*chosenParam).Consumed = false
			}
			continue
		}

		for _, chosenParam := range chosenParams {
			(*(*chosenParam).Payload).Status = types.CONSUMED
			messages = append(messages, (*(*chosenParam).Payload).Msg)

			if (*chosenParam).Payload.Ch != nil {
				syncPorts = append(syncPorts, (*(*chosenParam).Payload).Ch)
			}
		}

		retry = false
	}

	return messages, syncPorts, true
}

func findPending(params []*types.WrappedPayload) *types.WrappedPayload {
	for _, param := range params {
		if (*param).Consumed {
			continue
		}

		if (*(*param).Payload).Status == types.PENDING {
			(*param).Consumed = true
			return param
		}
	}

	return nil
}
