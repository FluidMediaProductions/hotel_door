package main

import (
	"net"
	"net/http"
	"github.com/golang/protobuf/proto"
	"github.com/fluidmediaproductions/hotel_door"
	"log"
	"time"
	"io"
	"io/ioutil"
	"bytes"
	"errors"
	"fmt"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"crypto/rand"
	"encoding/asn1"
	"os"
	"crypto"
	"encoding/hex"
)

const SecretTimeout = time.Minute

type Status struct {
	DoorNumber uint
	CurrentSecret []byte
	SecretGenTime time.Time
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
}

var status = &Status{}

type ActionHandler func([]byte) error

type Action struct {
	action door_comms.DoorAction
	handler ActionHandler
}

var actions = []Action{
	{
		action: door_comms.DoorAction_DOOR_UNLOCK,
		handler: unlockDoor,
	},
}

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

	macs, err := getMacAddr()
	if err != nil {
		return nil, err
	}

	wrappedMsg := &door_comms.ProtoMsg{
		Type: &msgType,
		Msg: msgBytes,
		Mac: proto.String(macs[0]),
		Sig: sig,
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

		status.DoorNumber = uint(*respMsg.DoorNum)

		if respMsg.GetActionRequired() {
			log.Println("Action required, getting action")
			actionId, actionType, actionData, err := getAction()
			if err != nil {
				log.Printf("Error getting action: %v\n", err)
			}
			err = handleAction(actionId, actionType, actionData)
			if err != nil {
				log.Printf("Error executing action: %v\n", err)
			}
		}
	}
}

func handleAction(actionId int64, actionType door_comms.DoorAction, actionData []byte) error {
	found := false
	for _, hander := range actions {
		if hander.action == actionType {
			err := hander.handler(actionData)
			if err != nil {
				return err
			} else {
				respMsg := &door_comms.ActionComplete{
					ActionId: proto.Int64(actionId),
				}
				_, err = sendMsg(respMsg, door_comms.MsgType_ACTION_COMPLETE, door_comms.MsgType_ACTION_COMPLETE_RESP)
				if err != nil {
					return err
				}
			}
			found = true
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("No action handler for %s\n", door_comms.DoorAction_name[int32(actionType)]))
	}
	return nil
}

func getAction() (int64, door_comms.DoorAction, []byte, error) {
	ping := &door_comms.GetAction{}

	resp, err := sendMsg(ping, door_comms.MsgType_GET_ACTION, door_comms.MsgType_GET_ACTION_RESP)

	respMsg := &door_comms.GetActionResp{}
	err = proto.Unmarshal(resp, respMsg)
	if err != nil {
		return 0, 0, nil, err
	}

	if respMsg.ActionType != nil {
		log.Printf("Action %s required\n", door_comms.DoorAction_name[int32(respMsg.GetActionType())])
		return respMsg.GetActionId(), respMsg.GetActionType(), respMsg.GetActionPayload(), nil
	}

	return 0, 0, nil, errors.New("no action to complete")
}

func unlockDoor(data []byte) error {
	if time.Since(status.SecretGenTime) > SecretTimeout {
		return errors.New("secret timed out")
	}
	msg := &door_comms.DoorUnlockAction{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return err
	}

	if bytes.Equal(msg.GetSecret(), status.CurrentSecret) {
		log.Println("Secret matches and si not out of date, unlocking door and expiring secret")
		status.SecretGenTime = time.Unix(0, 0)
		return nil
	} else {
		return errors.New("invalid secret")
	}

	return nil
}

func getKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if _, err := os.Stat("private.pem"); err == nil {
		key, err := loadPEMKey("private.pem")
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	} else {
		reader := rand.Reader
		bitSize := 2048

		key, err := rsa.GenerateKey(reader, bitSize)
		if err != nil {
			return nil, nil, err
		}

		publicKey := key.PublicKey

		err = savePEMKey("private.pem", key)
		if err != nil {
			return nil, nil, err
		}
		err = savePublicPEMKey("public.pem", publicKey)
		if err != nil {
			return nil, nil, err
		}
		return key, &publicKey, nil
	}
}

func loadPEMKey(fileName string) (*rsa.PrivateKey, error){
	outFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(outFile)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func savePEMKey(fileName string, key *rsa.PrivateKey) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		return err
	}
	return nil
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) error {
	asn1Bytes, err := asn1.Marshal(pubkey)
	if err != nil {
		return err
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	if err != nil {
		return err
	}
	return nil
}

func updateSecret() {
	ticker := time.NewTicker(time.Second * 60)
	for ; true; <- ticker.C {
		token := make([]byte, 32)
		rand.Read(token)
		status.CurrentSecret = token
		status.SecretGenTime = time.Now()
		log.Printf("New secret generated: %s\n", hex.EncodeToString(token))
	}
}

func main() {
	priv, pub, err := getKeys()
	if err != nil {
		log.Fatalf("Can't get encryption keys: %v\n", err)
	}
	status.PublicKey = pub
	status.PrivateKey = priv

	go updateSecret()
	pingServer()
}
