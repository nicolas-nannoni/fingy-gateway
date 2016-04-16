package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
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
		log.Errorf("Failed to set websocket upgrade: %+v", err)
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
		log.Error(err)
		connection.Close()
		return
	}

	go connection.ReadLoop()
	go connection.WriteLoop()
}
