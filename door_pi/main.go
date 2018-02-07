package main

import (
	"net"
	"net/http"
	"github.com/golang/protobuf/proto"
	"github.com/fluidmediaproductions/hotel_door"
	"log"
	"time"
	"bytes"
	"io"
	"io/ioutil"
)

type Status struct {
	DoorNumber uint
}

var status = &Status{}

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func sendMsg(msg proto.Message, msgType door_comms.MsgType) (*http.Response, error) {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	wrappedMsg := &door_comms.ProtoMsg{
		Type: &msgType,
		Msg: msgBytes,
	}

	wrappedMsgBytes, err := proto.Marshal(wrappedMsg)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:8000/proto", "application/x-google-protbuf", bytes.NewReader(wrappedMsgBytes))
	return resp, err
}

func readMsg(reader io.ReadCloser) ([]byte, *door_comms.MsgType, error) {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}
	defer reader.Close()

	newMsg := &door_comms.ProtoMsg{}
	err = proto.Unmarshal(body, newMsg)
	if err != nil {
		return nil, nil, err
	}

	return newMsg.Msg, newMsg.Type, nil
}

func pingServer() {
	ticker := time.NewTicker(time.Second * 1)
	for range ticker.C {
		macs, err := getMacAddr()
		if err != nil {
			log.Println(err)
			continue
		}
		ping := &door_comms.DoorPing{
			Mac:       proto.String(macs[0]),
			Timestamp: proto.Int64(time.Now().Unix()),
		}

		resp, err := sendMsg(ping, door_comms.MsgType_DOOR_PING)
		if err != nil {
			log.Printf("Error pinging server: %v\n", err)
			continue
		}
		msg, msgType, err := readMsg(resp.Body)
		if err != nil {
			log.Printf("Error pinging server: %v\n", err)
			continue
		}
		if *msgType != door_comms.MsgType_DOOR_PING_RESP {
			log.Printf("Wrong ping response type from server: %v\n", err)
			continue
		}

		respMsg := &door_comms.DoorPingResp{}
		err = proto.Unmarshal(msg, respMsg)
		if err != nil {
			log.Printf("Invalid ping response from server: %v\n", err)
			continue
		}

		if !*respMsg.Success {
			log.Printf("Error registering with server: %v\n", *respMsg.Error)
			continue
		}

		status.DoorNumber = uint(*respMsg.DoorNum)
	}
}

func main() {
	go pingServer()

	for {

	}
}
