package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/fluidmediaproductions/hotel_door"
	"crypto/rsa"
	"crypto/rand"
	"crypto"
	"net/http"
)

func sendMsg(msg proto.Message, msgType door_comms.MsgType, w http.ResponseWriter) error {
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
		Msg: msgBytes,
		Mac: proto.String(macs[0]),
		Sig: sig,
	}

	wrappedMsgBytes, err := proto.Marshal(wrappedMsg)
	if err != nil {
		return err
	}

	w.Write(wrappedMsgBytes)
	return nil
}
