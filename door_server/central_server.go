package main

import (
	"github.com/golang/protobuf/proto"
	"log"
	"time"
	"github.com/fluidmediaproductions/central_hotel_door_server/hotel_comms"
)

func connectToCentralServer() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		ping := &hotel_comms.HotelPing{
			Timestamp: proto.Int64(time.Now().Unix()),
		}

		resp, err := sendMsg(ping, hotel_comms.MsgType_HOTEL_PING, hotel_comms.MsgType_HOTEL_PING_RESP)
		if err != nil {
			log.Printf("Error pinging server: %v\n", err)
			continue
		}

		pingResp := &hotel_comms.HotelPingResp{}
		err = proto.Unmarshal(resp, pingResp)
		if err != nil {
			log.Printf("Error pinging server: %v\n", err)
			continue
		}

		if pingResp.GetActionRequired() {
			log.Println("Action required from central server")
			actions, _ := getActions()

			log.Println(actions)
		}
	}
}

func getActions() ([]*hotel_comms.Action, error) {
	ping := &hotel_comms.GetActions{}

	resp, err := sendMsg(ping, hotel_comms.MsgType_GET_ACTIONS, hotel_comms.MsgType_GET_ACTIONS_RESP)
	if err != nil {
		return nil, err
	}

	actionsResp := &hotel_comms.GetActionsResp{}
	err = proto.Unmarshal(resp, actionsResp)
	if err != nil {
		return nil, err
	}

	return actionsResp.GetActions(), nil
}
