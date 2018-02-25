package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"github.com/fluidmediaproductions/hotel_door"
	"errors"
	"crypto/rsa"
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

	action := &Action{
		PiID: pi.ID,
		Complete: proto.Bool(false),
	}
	var actions []*Action
	var actionCount int
	db.Where(&action).Find(&actions).Count(&actionCount)

	resp := &door_comms.DoorPingResp{
		Success: proto.Bool(true),
		DoorNum: proto.Uint32(door.Number),
		ActionRequired: proto.Bool(actionCount > 0),
	}

	w.WriteHeader(http.StatusOK)
	return sendMsg(resp, door_comms.MsgType_DOOR_PING_RESP, w)
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

	db.AutoMigrate(&Pi{}, &Door{}, &Action{})

	priv, pub, err := door_comms.GetKeys()
	if err != nil {
		log.Fatalf("Can't get encryption keys: %v\n", err)
	}
	status.PublicKey = pub
	status.PrivateKey = priv

	go checkPis()

	go serveStatic(":3001", "static/build")

	r := mux.NewRouter()
	r.Methods("POST").Path("/proto").HandlerFunc(protoServ)
	http.ListenAndServe(":8000", r)
}
