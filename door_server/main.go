package main

import (
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
)

var db *gorm.DB

type Pi struct {
	gorm.Model
	Mac string
	LastSeen time.Time
}

type Door struct {
	Pi Pi
	PiID uint
	Number uint32
}

type ProtoHandlerFunc func(msg []byte, w http.ResponseWriter) error

type ProtoHandler struct {
	msgType door_comms.MsgType
	handler ProtoHandlerFunc
}

var protoHandlers = []ProtoHandler{
	{
		msgType: door_comms.MsgType_DOOR_PING,
		handler: doorPing,
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
		if hander.msgType == *newMsg.Type {
			hander.handler(newMsg.Msg, w)
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

	wrappedMsg := &door_comms.ProtoMsg{
		Type: &msgType,
		Msg: msgBytes,
	}

	wrappedMsgBytes, err := proto.Marshal(wrappedMsg)
	if err != nil {
		return err
	}

	w.Write(wrappedMsgBytes)
	return nil
}

func doorPing(msg []byte, w http.ResponseWriter) error {
	newMsg := &door_comms.DoorPing{}
	err := proto.Unmarshal(msg, newMsg)
	if err != nil {
		return err
	}

	if time.Since(time.Unix(*newMsg.Timestamp, 0)) > time.Hour {
		log.Printf("Pi %v out of sync with server time\n", newMsg.Mac)

		resp := &door_comms.DoorPingResp{
			Success: proto.Bool(false),
			Error: proto.String("time out of sync"),
		}
		w.WriteHeader(http.StatusNotAcceptable)
		sendMsg(resp, door_comms.MsgType_DOOR_PING_RESP, w)
		return errors.New("pi out of sync")
	}

	pi := &Pi{
		Mac: *newMsg.Mac,
	}


	db.Where(pi).First(pi)
	pi.LastSeen = time.Now()
	db.Save(pi)

	door := &Door{
		PiID: pi.ID,
	}
	db.First(&door)

	resp := &door_comms.DoorPingResp{
		Success: proto.Bool(true),
		DoorNum: proto.Uint32(door.Number),
	}

	w.WriteHeader(http.StatusOK)
	return sendMsg(resp, door_comms.MsgType_DOOR_PING_RESP, w)
}

func checkPis() {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		pis := make([]*Pi, 0)
		db.Find(&pis)

		log.Println("*** Current pis ***")
		if len(pis) == 0 {
			log.Println("None")
		} else {
			for _, pi := range pis {
				log.Printf("MAC: %v Last seen: %v\n", pi.Mac, pi.LastSeen)
				if time.Since(pi.LastSeen) > time.Minute {
					log.Printf("Removing pi %v, too old\n", pi.Mac)
					db.Delete(&pi)
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

	db.AutoMigrate(&Pi{}, &Door{})

	go checkPis()

	r := mux.NewRouter()
	r.Methods("POST").Path("/proto").HandlerFunc(protoServ)
	http.ListenAndServe(":8000", r)
}
