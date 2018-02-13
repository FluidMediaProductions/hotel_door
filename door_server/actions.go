package main

import (
	"github.com/fluidmediaproductions/hotel_door"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Action struct {
	gorm.Model
	Type int
	Payload []byte
	Pi Pi
	PiID uint
	Complete *bool
	Success *bool
}

func getAction(pi *Pi, msg []byte, _ []byte, w http.ResponseWriter) error {
	newMsg := &door_comms.GetAction{}
	err := proto.Unmarshal(msg, newMsg)
	if err != nil {
		return err
	}

	action := &Action{
		PiID: pi.ID,
		Complete: proto.Bool(false),
	}
	var actions []*Action
	var actionCount int
	db.Where(action).Find(&actions).Count(&actionCount)

	var resp *door_comms.GetActionResp
	if actionCount < 1 {
		resp = &door_comms.GetActionResp{}
	} else {
		actionType := door_comms.DoorAction(actions[0].Type)
		resp = &door_comms.GetActionResp{
			ActionType: &actionType,
			ActionPayload: actions[0].Payload,
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
	db.First(action, newMsg.ActionId).Count(&actionCount)

	var resp *door_comms.ActionCompleteResp
	if actionCount < 1 {
		resp = &door_comms.ActionCompleteResp{}
	} else {
		action.Complete = proto.Bool(true)
		action.Success = proto.Bool(false)
		db.Save(action)
		resp = &door_comms.ActionCompleteResp{}
	}

	w.WriteHeader(http.StatusOK)
	return sendMsg(resp, door_comms.MsgType_ACTION_COMPLETE_RESP, w)
}
