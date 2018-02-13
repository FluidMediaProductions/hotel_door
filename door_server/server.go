package main

import (
	"github.com/fluidmediaproductions/hotel_door"
	"crypto/x509"
	"crypto/rsa"
	"github.com/golang/protobuf/proto"
	"crypto"
	"errors"
	"net/http"
	"log"
	"io/ioutil"
)

type ProtoHandlerFunc func(pi *Pi, msg []byte, sig []byte, w http.ResponseWriter) error

type ProtoHandler struct {
	msgType door_comms.MsgType
	handler ProtoHandlerFunc
	noSigRequired bool
}

var protoHandlers = []ProtoHandler{
	{
		msgType: door_comms.MsgType_DOOR_PING,
		handler: doorPing,
		noSigRequired: true,
	},
	{
		msgType: door_comms.MsgType_GET_ACTION,
		handler: getAction,
	},
	{
		msgType: door_comms.MsgType_ACTION_COMPLETE,
		handler: actionComplete,
	},
}

func protoServ(w http.ResponseWriter, r *http.Request) {
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("body error:", readErr)
		return
	}
	defer r.Body.Close()

	newMsg := &door_comms.ProtoMsg{}
	err := proto.Unmarshal(body, newMsg)
	if err != nil {
		log.Println("Proto error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, hander := range protoHandlers {
		if hander.msgType == newMsg.GetType() {
			pi := &Pi{
				Mac: *newMsg.Mac,
			}
			db.First(pi)

			if !hander.noSigRequired {
				err := verifySignature(newMsg.Msg, newMsg.Sig, pi.PublicKey)
				if err != nil {
					log.Printf("Unable to verify signature from %s: %v\n", pi.Mac, err)
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
			err := hander.handler(pi, newMsg.GetMsg(), newMsg.GetSig(), w)
			if err != nil {
				log.Printf("Error on handler for %s: %v\n", door_comms.MsgType_name[int32(newMsg.GetType())], err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func verifySignature(msg []byte, sig []byte, pubKey []byte) error {
	pub, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		return err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		hash := crypto.SHA256
		h := hash.New()
		h.Write(msg)
		hashed := h.Sum(nil)
		err := rsa.VerifyPKCS1v15(pub, hash, hashed, sig)

		if err != nil {
			return err
		}
	default:
		return errors.New("invalid public key type")
	}
	return nil
}

