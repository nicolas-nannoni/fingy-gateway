package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nicolas-nannoni/fingy-gateway/events"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

const (
	registerBufferSize    = 100
	unregisterBufferSize  = 100
	messageSendBufferSize = 512
)

type connection struct {
	id        uuid.UUID
	deviceId  string
	serviceId string
	ws        *websocket.Conn
	send      chan []byte
}

func (c *connection) String() string {
	return fmt.Sprintf("{id: %s, deviceId: %s}", c.id, c.deviceId)
}

func (c *connection) Close() {
	c.write(websocket.CloseMessage, []byte{})
	close(c.send)
}

func (c *connection) WriteLoop() {
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

func (c *connection) ReadLoop() {

	c.ws.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error in connection %s: %v", c, err)
				Registry.unregisterConnection(c)
			}
			break
		}
		c.dispatchReceivedMessage(message)
	}
	log.Printf("Read loop closed %s", c)
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) dispatchReceivedMessage(msg []byte) {

	log.Printf("Received message: %s on connection %s", msg, c)
	var evt events.Event
	err := json.Unmarshal(msg, &evt)
	if err != nil {
		log.Print(err)
		return
	}

	resp, err := Registry.Dispatch(c.serviceId, c.deviceId, &evt)
	log.Printf("Response: %v, error: %v", resp, err)
}
