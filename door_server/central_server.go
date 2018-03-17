package main

import (
	"crypto/x509"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
	"github.com/fluidmediaproductions/central_hotel_door_server/hotel_comms"
)

func connectToCentralServer() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		pub, err := x509.MarshalPKIXPublicKey(status.PublicKey)
		if err != nil {
			log.Printf("Cannot marshal public key: %v\n", err)
			continue
		}
		ping := &hotel_comms.HotelPing{
			Timestamp: proto.Int64(time.Now().Unix()),
			PublicKey: pub,
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

		log.Println(pingResp)
	}
}
