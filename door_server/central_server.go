package main

import (
	"github.com/golang/protobuf/proto"
	"log"
	"time"
	"log"
	"github.com/fluidmediaproductions/central_hotel_door_server/hotel_comms"
	"github.com/golang/protobuf/proto"
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
			log.Printf("Error updating doors server: %v\n", err)
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

	log.Println(pingResp)
	return nil
}

func updateDoors() error {
	ping := &hotel_comms.GetDoors{}

	resp, err := sendMsg(ping, hotel_comms.MsgType_GET_DOORS, hotel_comms.MsgType_GET_DOORS_RESP)
	if err != nil {
		return err
	}

	pingResp := &hotel_comms.GetActionsResp{}
	err = proto.Unmarshal(resp, pingResp)
	if err != nil {
		return err
	}

	log.Println(pingResp)
	return nil
}
