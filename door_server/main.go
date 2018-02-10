package main

import (
	"net"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"log"
	"github.com/fluidmediaproductions/hotel_door"
	"errors"
	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
	"crypto"
	"encoding/pem"
	"os"
	"encoding/asn1"
)

var db *gorm.DB

type Status struct {
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
}

var status = &Status{}

type Pi struct {
	gorm.Model
	Mac string
	LastSeen time.Time
	Online bool
	PublicKey []byte
}

type Door struct {
	gorm.Model
	Pi Pi
	PiID uint
	Number uint32
}

type PendingAction struct {
	gorm.Model
	Type int
	Payload []byte
	Pi Pi
	PiID uint
}

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

func sendMsg(msg proto.Message, msgType door_comms.MsgType, w http.ResponseWriter) error {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	macs, err := getMacAddr()
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

func doorPing(pi *Pi, msg []byte, sig []byte, w http.ResponseWriter) error {
	newMsg := &door_comms.DoorPing{}
	err := proto.Unmarshal(msg, newMsg)
	if err != nil {
		return err
	}

	err = verifySignature(msg, sig, newMsg.GetPublicKey())
	if err != nil {
		log.Printf("Unable to verify pi %v signature\n", pi.Mac)

		resp := &door_comms.DoorPingResp{
			Success: proto.Bool(false),
			Error: proto.String("invalid signature"),
		}
		w.WriteHeader(http.StatusNotAcceptable)
		sendMsg(resp, door_comms.MsgType_DOOR_PING_RESP, w)
		return err
	}

	if time.Since(time.Unix(*newMsg.Timestamp, 0)) > time.Hour {
		log.Printf("Pi %v out of sync with server time\n", pi.Mac)

		resp := &door_comms.DoorPingResp{
			Success: proto.Bool(false),
			Error: proto.String("time out of sync"),
		}
		w.WriteHeader(http.StatusNotAcceptable)
		sendMsg(resp, door_comms.MsgType_DOOR_PING_RESP, w)
		return errors.New("pi out of sync")
	}

	pi.LastSeen = time.Now()
	pi.PublicKey = newMsg.GetPublicKey()
	db.Save(pi)

	door := &Door{
		PiID: pi.ID,
	}
	db.First(&door)

	action := &PendingAction{
		PiID: pi.ID,
	}
	var actions []*PendingAction
	var actionCount int
	db.Find(&actions, action).Count(&actionCount)

	resp := &door_comms.DoorPingResp{
		Success: proto.Bool(true),
		DoorNum: proto.Uint32(door.Number),
		ActionRequired: proto.Bool(actionCount > 0),
	}

	w.WriteHeader(http.StatusOK)
	return sendMsg(resp, door_comms.MsgType_DOOR_PING_RESP, w)
}

func getAction(pi *Pi, msg []byte, _ []byte, w http.ResponseWriter) error {
	newMsg := &door_comms.GetAction{}
	err := proto.Unmarshal(msg, newMsg)
	if err != nil {
		return err
	}

	action := &PendingAction{
		PiID: pi.ID,
	}
	var actions []*PendingAction
	var actionCount int
	db.Find(&actions, action).Count(&actionCount)

	var resp *door_comms.GetActionResp
	if actionCount < 1 {
		resp = &door_comms.GetActionResp{}
	} else {
		actionType := door_comms.DoorAction(actions[0].Type)
		resp = &door_comms.GetActionResp{
			ActionType: &actionType,
			ActionPayload: actions[0].Payload,
		}
	}

	log.Println(resp)

	w.WriteHeader(http.StatusOK)
	return sendMsg(resp, door_comms.MsgType_GET_ACTION_RESP, w)
}

func actionComplete(pi *Pi, msg []byte, _ []byte, w http.ResponseWriter) error {
	newMsg := &door_comms.ActionComplete{}
	err := proto.Unmarshal(msg, newMsg)
	if err != nil {
		return err
	}

	action := &PendingAction{
		PiID: pi.ID,
	}
	var actionCount int
	db.First(action, newMsg.ActionId).Count(&actionCount)

	var resp *door_comms.ActionCompleteResp
	if actionCount < 1 {
		resp = &door_comms.ActionCompleteResp{}
	} else {
		db.Delete(action)
		resp = &door_comms.ActionCompleteResp{}
	}

	log.Println(resp)

	w.WriteHeader(http.StatusOK)
	return sendMsg(resp, door_comms.MsgType_ACTION_COMPLETE_RESP, w)
}

func checkPis() {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		pis := make([]*Pi, 0)
		db.Find(&pis, &Pi{Online: true})

		log.Println("*** Current pis ***")
		if len(pis) == 0 {
			log.Println("None")
		} else {
			for _, pi := range pis {
				log.Printf("MAC: %v Last seen: %v\n", pi.Mac, pi.LastSeen)
				if time.Since(pi.LastSeen) > time.Minute {
					log.Printf("Removing pi %v, too old\n", pi.Mac)
					pi.Online = false
					db.Save(&pi)
				}
			}
		}
	}
}

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Pi{}, &Door{}, &PendingAction{})

	priv, pub, err := getKeys()
	if err != nil {
		log.Fatalf("Can't get encryption keys: %v\n", err)
	}
	status.PublicKey = pub
	status.PrivateKey = priv

	go checkPis()

	r := mux.NewRouter()
	r.Methods("POST").Path("/proto").HandlerFunc(protoServ)
	http.ListenAndServe(":8000", r)
}
