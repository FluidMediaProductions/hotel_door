package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/fluidmediaproductions/central_hotel_door_server/hotel_comms"
	"github.com/fluidmediaproductions/hotel_door"
)

const CentralServer = "http://localhost:8081"

func sendMsgResp(msg proto.Message, msgType door_comms.MsgType, w http.ResponseWriter) error {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	macs, err := door_comms.GetMacAddr()
	if err != nil {
		return err
	}

	reader := rand.Reader
	hash := crypto.SHA256
	h := hash.New()
	h.Write(msgBytes)
	hashed := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(reader, status.PrivateKey, hash, hashed)
	if err != nil {
		return err
	}

	wrappedMsg := &door_comms.ProtoMsg{
		Type: &msgType,
		Msg:  msgBytes,
		Mac:  proto.String(macs[0]),
		Sig:  sig,
	}

	wrappedMsgBytes, err := proto.Marshal(wrappedMsg)
	if err != nil {
		return err
	}

	w.Write(wrappedMsgBytes)
	return nil
}

func sendMsg(msg proto.Message, msgType hotel_comms.MsgType, respMsgType hotel_comms.MsgType) ([]byte, error) {
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

	wrappedMsg := &hotel_comms.ProtoMsg{
		Type: &msgType,
		Msg:  msgBytes,
		UUID:  proto.String(status.UUID),
		Sig:  sig,
	}

	wrappedMsgBytes, err := proto.Marshal(wrappedMsg)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(CentralServer+"/proto", "application/x-google-protobuf", bytes.NewReader(wrappedMsgBytes))

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

func readMsg(reader io.ReadCloser) ([]byte, hotel_comms.MsgType, error) {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, 0, err
	}
	defer reader.Close()

	newMsg := &hotel_comms.ProtoMsg{}
	err = proto.Unmarshal(body, newMsg)
	if err != nil {
		return nil, 0, err
	}

	return newMsg.GetMsg(), newMsg.GetType(), nil
}
