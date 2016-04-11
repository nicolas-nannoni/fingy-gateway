package main

import (
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{}

func socketHandler(serviceId string, deviceId string, w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	connection := connection{
		id:        uuid.NewV1(),
		deviceId:  deviceId,
		serviceId: serviceId,
		ws:        ws,
		send:      make(chan []byte),
	}

	err = Registry.registerConnection(&connection)
	if err != nil {
		log.Print(err)
		connection.Close()
		return
	}

	go connection.ReadLoop()
	go connection.WriteLoop()
}
