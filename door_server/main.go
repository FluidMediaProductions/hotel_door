package main

import (
	"crypto/rsa"
	"errors"
	"github.com/fluidmediaproductions/hotel_door"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	uuid2 "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
	"bytes"
	"github.com/dgrijalva/jwt-go"
	"encoding/base64"
)

const JWTSecret = "RSQikzffdBsJjjtrzIHbrxI6UD1+BgZgOBGY7H8O2BkOsFsES1s5zStR1Qn6mseswRTTbT+sdwKLk5jFSpkQtQ=="
var JWTSecretBytes []byte

type JWTClaims struct {
	User *User        `json:"user"`
	jwt.StandardClaims
}

var db *gorm.DB

type Status struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	UUID       string
}

var status = &Status{}

type Config struct {
	gorm.Model
	Name  string
	Value string
}

type Pi struct {
	gorm.Model
	Mac       string    `json:"mac"`
	LastSeen  time.Time `json:"lastSeen"`
	Online    bool      `json:"online"`
	PublicKey []byte    `json:"-"`
}

type Door struct {
	gorm.Model
	Pi   *Pi    `json:"pi"`
	PiID uint   `json:"piId"`
	Name string `json:"name"`
}

type Action struct {
	gorm.Model
	Pi       *Pi    `json:"pi"`
	PiID     uint   `json:"piId"`
	Type     int    `json:"type"`
	Payload  []byte `json:"payload"`
	Complete bool   `json:"complete"`
	Success  bool   `json:"success"`
}

type User struct {
	gorm.Model
	User string        `json:"user"`
	Pass string        `json:"-"`
	Name string        `json:"name"`
	IsAdmin bool       `json:"isAdmin"`
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
			Error:   proto.String("invalid signature"),
		}
		w.WriteHeader(http.StatusNotAcceptable)
		sendMsgResp(resp, door_comms.MsgType_DOOR_PING_RESP, w)
		return err
	}

	if time.Since(time.Unix(*newMsg.Timestamp, 0)) > time.Hour {
		log.Printf("Pi %v out of sync with server time\n", pi.Mac)

		resp := &door_comms.DoorPingResp{
			Success: proto.Bool(false),
			Error:   proto.String("time out of sync"),
		}
		w.WriteHeader(http.StatusNotAcceptable)
		sendMsgResp(resp, door_comms.MsgType_DOOR_PING_RESP, w)
		return errors.New("pi out of sync")
	}

	if pi.PublicKey != nil {
		if !bytes.Equal(newMsg.GetPublicKey(), pi.PublicKey) {
			log.Printf("Pi %v already registered with different public key\n", pi.Mac)

			resp := &door_comms.DoorPingResp{
				Success: proto.Bool(false),
				Error: proto.String("already registered"),
			}
			w.WriteHeader(http.StatusForbidden)
			sendMsgResp(resp, door_comms.MsgType_DOOR_PING_RESP, w)
			return err
		}
	}

	pi.LastSeen = time.Now()
	pi.PublicKey = newMsg.GetPublicKey()
	pi.Online = true
	db.Save(pi)

	door := &Door{
		PiID: pi.ID,
	}
	db.First(&door, &door)

	action := &Action{}
	var actionCount int
	db.Where(map[string]interface{}{"pi_id": pi.ID, "complete": false}).Find(&action).Count(&actionCount)

	resp := &door_comms.DoorPingResp{
		Success:        proto.Bool(true),
		DoorName:        proto.String(door.Name),
		ActionRequired: proto.Bool(actionCount > 0),
	}

	w.WriteHeader(http.StatusOK)
	return sendMsgResp(resp, door_comms.MsgType_DOOR_PING_RESP, w)
}

func checkPis() {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		pis := make([]*Pi, 0)
		db.Find(&pis, &Pi{Online: true})

		for _, pi := range pis {
			if time.Since(pi.LastSeen) > time.Minute {
				log.Printf("Removing pi %v, too old\n", pi.Mac)
				pi.Online = false
				db.Save(&pi)
			}
		}
	}
}

func getUUID() (string, error) {
	config := &Config{}
	err := db.First(config, &Config{
		Name: "UUID",
	}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uuid, err := uuid2.NewV4()
			if err != nil {
				return "", err
			}
			config = &Config{
				Name:  "UUID",
				Value: uuid.String(),
			}
			db.Save(config)
		} else {
			return "", err
		}
	}
	return config.Value, nil
}

func main() {
	var err error

	JWTSecretBytes, err = base64.StdEncoding.DecodeString(JWTSecret)
	if err != nil {
		panic(err)
	}

	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Pi{}, &Door{}, &Action{}, &User{}, &Config{})

	priv, pub, err := door_comms.GetKeys()
	if err != nil {
		log.Fatalf("Can't get encryption keys: %v\n", err)
	}
	status.PublicKey = pub
	status.PrivateKey = priv

	uuid, err := getUUID()
	if err != nil {
		log.Fatalf("Can't get UUID: %v\n", err)
	}
	status.UUID = uuid

	go connectToCentralServer()
	go checkPis()

	go serveStatic(":3001", "static/build")

	r := mux.NewRouter()
	r.Methods("POST").Path("/proto").HandlerFunc(protoServ)
	log.Println("Listening on :8000")
	log.Fatalln(http.ListenAndServe(":8000", r))
}
