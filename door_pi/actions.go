package main

import (
	"errors"
	"fmt"
	"github.com/fluidmediaproductions/hotel_door"
	"github.com/golang/protobuf/proto"
	"log"
)

type ActionHandler func([]byte) error

type Action struct {
	action  door_comms.DoorAction
	handler ActionHandler
}

var actions = []Action{
	{
		action:  door_comms.DoorAction_DOOR_UNLOCK,
		handler: unlockDoor,
	},
}

func handleAction(actionId int64, actionType door_comms.DoorAction, actionData []byte) error {
	found := false
	for _, hander := range actions {
		if hander.action == actionType {
			err := hander.handler(actionData)
			var respMsg *door_comms.ActionComplete
			var success bool
			if err != nil {
				success = false
				respMsg = &door_comms.ActionComplete{
					ActionId: proto.Int64(actionId),
					Success:  proto.Bool(false),
				}
			} else {
				success = true
				respMsg = &door_comms.ActionComplete{
					ActionId: proto.Int64(actionId),
					Success:  proto.Bool(true),
				}
			}
			_, err2 := sendMsg(respMsg, door_comms.MsgType_ACTION_COMPLETE, door_comms.MsgType_ACTION_COMPLETE_RESP)
			if !success {
				return err
			}
			if err2 != nil {
				return err2
			}
			found = true
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("no action handler for %s", door_comms.DoorAction_name[int32(actionType)]))
	}
	return nil
}

func getAction() (int64, door_comms.DoorAction, []byte, error) {
	ping := &door_comms.GetAction{}

	resp, err := sendMsg(ping, door_comms.MsgType_GET_ACTION, door_comms.MsgType_GET_ACTION_RESP)

	respMsg := &door_comms.GetActionResp{}
	err = proto.Unmarshal(resp, respMsg)
	if err != nil {
		return 0, 0, nil, err
	}

	if respMsg.ActionType != nil {
		log.Printf("Action %s required\n", door_comms.DoorAction_name[int32(respMsg.GetActionType())])
		return respMsg.GetActionId(), respMsg.GetActionType(), respMsg.GetActionPayload(), nil
	}

	return 0, 0, nil, errors.New("no action to complete")
}
