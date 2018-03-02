package main

import (
	"log"
)

func unlockDoor(_ []byte) error {
	log.Println("Unlocking door")

	status.gui.SetDoorOpening()

	return nil
}
