package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"github.com/fluidmediaproductions/hotel_door"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"net/http"
)

func sendMsg(msg proto.Message, msgType door_comms.MsgType, respMsgType door_comms.MsgType) ([]byte, error) {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	reader := rand.Reader
	hash := crypto.SHA256
	h := hash.New()
	h.Write(msgBytes)
	hashed := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(reader, status.PrivateKey, hash, hashed)
	if err != nil {
		return nil, err
	}

	macs, err := door_comms.GetMacAddr()
	if err != nil {
		return nil, err
	}

	wrappedMsg := &door_comms.ProtoMsg{
		Type: &msgType,
		Msg:  msgBytes,
		Mac:  proto.String(macs[0]),
		Sig:  sig,
	}

	wrappedMsgBytes, err := proto.Marshal(wrappedMsg)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:8000/proto", "application/x-google-protobuf", bytes.NewReader(wrappedMsgBytes))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, respType, err := readMsg(resp.Body)
	if err != nil {
		return nil, err
	}
	if respType != respMsgType {
		return nil, errors.New("wrong response type")
	}
	return respBytes, nil
}

func readMsg(reader io.ReadCloser) ([]byte, door_comms.MsgType, error) {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, 0, err
	}
	defer reader.Close()

	newMsg := &door_comms.ProtoMsg{}
	err = proto.Unmarshal(body, newMsg)
	if err != nil {
		return nil, 0, err
	}

	return newMsg.GetMsg(), newMsg.GetType(), nil
}
