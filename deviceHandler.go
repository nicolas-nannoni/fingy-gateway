package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/nicolas-nannoni/fingy-server/events"
	"github.com/nicolas-nannoni/fingy-server/services"
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

func socketHandler(deviceId string, w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	connection := connection{id: uuid.NewV1(), deviceId: deviceId, ws: ws, send: make(chan []byte)}
	reg.register <- &connection

	go connection.readLoop()
	go connection.writeLoop()
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) close() {
	c.write(websocket.CloseMessage, []byte{})
	close(c.send)
}

func (c *connection) writeLoop() {
Loop:
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				break Loop
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				log.Fatalf("Error while sending message to connection %s", c)
				return
			}
		}
	}
	log.Printf("Write loop closed %s", c)
}

func (c *connection) readLoop() {

	c.ws.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error in connection %s: %v", c, err)
				reg.unregisterConnection(c)
			}
			break
		}
		c.dispatchReceivedMessage(message)
	}
	log.Printf("Read loop closed %s", c)
}

func (c *connection) dispatchReceivedMessage(msg []byte) {

	log.Printf("Received message: %s on connection %s", msg, c)
	var evt events.Event
	err := json.Unmarshal(msg, &evt)
	if err != nil {
		log.Print(err)
		return
	}

	resp, err := services.Registry.Dispatch(&evt)
	log.Printf("Response: %v, error: %v", resp, err)
}
