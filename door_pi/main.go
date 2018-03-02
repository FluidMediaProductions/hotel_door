package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/fluidmediaproductions/hotel_door"
	"log"
	"time"
	"crypto/rsa"
	"crypto/x509"
	"github.com/fluidmediaproductions/hotel_door/door_pi/gui"
)

type Status struct {
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
	gui *gui.GUI
}

var status = &Status{}

func pingServer() {
	ticker := time.NewTicker(time.Second * 1)
	for range ticker.C {
		pub, err := x509.MarshalPKIXPublicKey(status.PublicKey)
		if err != nil {
			log.Printf("Cannot marshal public key: %v\n", err)
			continue
		}
		ping := &door_comms.DoorPing{
			Timestamp: proto.Int64(time.Now().Unix()),
			PublicKey: pub,
		}

		resp, err := sendMsg(ping, door_comms.MsgType_DOOR_PING, door_comms.MsgType_DOOR_PING_RESP)
		if err != nil {
			log.Printf("Error pinging server: %v\n", err)
			continue
		}

		respMsg := &door_comms.DoorPingResp{}
		err = proto.Unmarshal(resp, respMsg)
		if err != nil {
			log.Printf("Invalid ping response from server: %v\n", err)
			continue
		}

		if !respMsg.GetSuccess() {
			log.Printf("Error registering with server: %v\n", *respMsg.Error)
			continue
		}

		status.gui.SetRoomNum(int(*respMsg.DoorNum))

		if respMsg.GetActionRequired() {
			log.Println("Action required, getting action")
			actionId, actionType, actionData, err := getAction()
			if err != nil {
				log.Printf("Error getting action: %v\n", err)
				continue
			}
			err = handleAction(actionId, actionType, actionData)
			if err != nil {
				log.Printf("Error executing action: %v\n", err)
				continue
			}
		}
	}
}

func main() {
	priv, pub, err := door_comms.GetKeys()
	if err != nil {
		log.Fatalf("Can't get encryption keys: %v\n", err)
	}
	status.PublicKey = pub
	status.PrivateKey = priv

	status.gui = gui.NewGui()

	go pingServer()
	status.gui.Start()
}
