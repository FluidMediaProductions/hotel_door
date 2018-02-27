package main

import (
	"github.com/fluidmediaproductions/hotel_door"
	"github.com/golang/protobuf/proto"
	"net/http"
)

func getAction(pi *Pi, msg []byte, _ []byte, w http.ResponseWriter) error {
	newMsg := &door_comms.GetAction{}
	err := proto.Unmarshal(msg, newMsg)
	if err != nil {
		return err
	}

	action := &Action{}
	var actionCount int
	db.Where(map[string]interface{}{"pi_id": pi.ID, "complete": false}).Find(&action).Count(&actionCount)

	var resp *door_comms.GetActionResp
	if actionCount < 1 {
		resp = &door_comms.GetActionResp{}
	} else {
		actionType := door_comms.DoorAction(action.Type)
		resp = &door_comms.GetActionResp{
			ActionId: proto.Int64(int64(action.ID)),
			ActionType: &actionType,
			ActionPayload: action.Payload,
		}
	}

	w.WriteHeader(http.StatusOK)
	return sendMsg(resp, door_comms.MsgType_GET_ACTION_RESP, w)
}

func actionComplete(pi *Pi, msg []byte, _ []byte, w http.ResponseWriter) error {
	newMsg := &door_comms.ActionComplete{}
	err := proto.Unmarshal(msg, newMsg)
	if err != nil {
		return err
	}

	action := &Action{
		PiID: pi.ID,
	}
	var actionCount int
	db.First(&action).Count(&actionCount)

	var resp *door_comms.ActionCompleteResp
	if actionCount < 1 {
		resp = &door_comms.ActionCompleteResp{}
	} else {
		action.Complete = true
		action.Success = *newMsg.Success
		db.Save(action)
		resp = &door_comms.ActionCompleteResp{}
	}

	w.WriteHeader(http.StatusOK)
	return sendMsg(resp, door_comms.MsgType_ACTION_COMPLETE_RESP, w)
}
