package main

import (
	"github.com/fluidmediaproductions/hotel_door"
	"github.com/golang/protobuf/proto"
	"bytes"
	"encoding/hex"
	"time"
	"errors"
	"log"
	"crypto/rand"
)

const SecretTimeout = time.Minute

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
