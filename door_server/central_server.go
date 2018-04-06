package main

import (
	"github.com/golang/protobuf/proto"
	"log"
	"time"
	"github.com/fluidmediaproductions/central_hotel_door_server/hotel_comms"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/fluidmediaproductions/hotel_door"
)

func connectToCentralServer() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		err := pingCentralServer()
		if err != nil {
			log.Printf("Error pinging server: %v\n", err)
			continue
		}

		err = updateDoors()
		if err != nil {
			log.Printf("Error updating doors from server: %v\n", err)
			continue
		}
	}
}

func pingCentralServer() error {
	ping := &hotel_comms.HotelPing{
		Timestamp: proto.Int64(time.Now().Unix()),
	}

	resp, err := sendMsg(ping, hotel_comms.MsgType_HOTEL_PING, hotel_comms.MsgType_HOTEL_PING_RESP)
	if err != nil {
		return err
	}

	pingResp := &hotel_comms.HotelPingResp{}
	err = proto.Unmarshal(resp, pingResp)
	if err != nil {
		return err
	}

	if pingResp.GetSuccess() != true {
		return errors.New(pingResp.GetError())
	} else {
		if pingResp.GetActionRequired() {
			err := getActionsFromCentralServer()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getActionsFromCentralServer() error {
	ping := &hotel_comms.GetActions{}

	resp, err := sendMsg(ping, hotel_comms.MsgType_GET_ACTIONS, hotel_comms.MsgType_GET_ACTIONS_RESP)
	if err != nil {
		return err
	}

	pingResp := &hotel_comms.GetActionsResp{}
	err = proto.Unmarshal(resp, pingResp)
	if err != nil {
		return err
	}

	for _, action := range pingResp.GetActions() {
		var err error
		if action.GetType() == hotel_comms.ActionType_ROOM_UNLOCK {
			err = unlockRoom(action.GetId())
		}
		if err == nil {
			msg := &hotel_comms.ActionComplete{
				ActionType: action.Type,
				ActionId: action.Id,
				Success: proto.Bool(true),
			}
			_, err := sendMsg(msg, hotel_comms.MsgType_ACTION_COMPLETE, hotel_comms.MsgType_ACTION_COMPLETE_RESP)
			if err != nil {
				return err
			}
		} else {
			msg := &hotel_comms.ActionComplete{
				ActionType: action.Type,
				ActionId: action.Id,
				Success: proto.Bool(false),
			}
			log.Println(err)
			_, err := sendMsg(msg, hotel_comms.MsgType_ACTION_COMPLETE, hotel_comms.MsgType_ACTION_COMPLETE_RESP)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func unlockRoom(roomId int64) error {
	door := &Door{}
	err := db.First(door, roomId).Error
	if err != nil {
		return err
	}

	action := &Action{
		PiID: door.PiID,
		Type: int(door_comms.DoorAction_DOOR_UNLOCK),
	}
	err = db.Create(action).Error
	if err != nil {
		return err
	}
	return nil
}

func updateDoors() error {
	ping := &hotel_comms.GetDoors{}

	resp, err := sendMsg(ping, hotel_comms.MsgType_GET_DOORS, hotel_comms.MsgType_GET_DOORS_RESP)
	if err != nil {
		return err
	}

	pingResp := &hotel_comms.GetDoorsResp{}
	err = proto.Unmarshal(resp, pingResp)
	if err != nil {
		return err
	}

	for _, door := range pingResp.GetDoors() {
		dbDoor := &Door{
			Model: gorm.Model{
				ID: uint(door.GetId()),
			},
		}

		err := db.First(dbDoor, door.GetId()).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				dbDoor.Name = door.GetName()
				err := db.Save(dbDoor).Error
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if dbDoor.Name != door.GetName() {
			dbDoor.Name = door.GetName()
			err := db.Save(dbDoor).Error
			if err != nil {
				return err
			}
		}
	}

	dbDoors := make([]*Door, 0)
	err = db.Find(&dbDoors).Error
	if err != nil {
		return err
	}
	for _, dbDoor := range dbDoors {
		doorExists := false
		for _, door := range pingResp.GetDoors() {
			if uint(door.GetId()) == dbDoor.ID {
				doorExists = true
				break
			}
		}
		if !doorExists {
			err := db.Delete(dbDoor).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
