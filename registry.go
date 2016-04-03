package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nicolas-nannoni/fingy-server/events"
	"github.com/satori/go.uuid"
	"log"
)

const (
	registerBufferSize    = 100
	unregisterBufferSize  = 100
	messageSendBufferSize = 512
)

type connection struct {
	id       uuid.UUID
	deviceId string
	ws       *websocket.Conn
	send     chan []byte
}

func (c *connection) String() string {
	return fmt.Sprintf("{id: %s, deviceId: %s}", c.id, c.deviceId)
}

type registry struct {
	connections map[string]*connection
	register    chan *connection
	unregister  chan *connection
}

var reg = registry{
	connections: make(map[string]*connection),
	register:    make(chan *connection, registerBufferSize),
	unregister:  make(chan *connection, unregisterBufferSize),
}

func (r *registry) run() {

	for {
		select {
		case c := <-r.register:
			r.registerConnection(c)
		case c := <-r.unregister:
			r.unregisterConnection(c)
		}
	}
}

func (r *registry) Send(deviceId string, evt *events.Event) (err error) {

	c, ok := r.connections[deviceId]
	if !ok {
		return fmt.Errorf("The device with id %s is not registered", deviceId)
	}

	evt.PrepareForSend()
	err = evt.Verify()
	if err != nil {
		return err
	}

	msg, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("The event %s could not be serialized to JSON: %v", evt, err)
	}
	log.Printf("Pushing message %s to send queue of %s", msg, c)
	c.send <- msg

	return
}

func (r *registry) registerConnection(c *connection) {

	log.Printf("Registering connection %s", c)
	if existingConn, ok := r.connections[c.deviceId]; ok {
		log.Printf("Existing registration for device %s. Closing old connection %s", c.deviceId, existingConn)
		r.unregisterConnection(existingConn)
	}

	r.connections[c.deviceId] = c
}

func (r *registry) unregisterConnection(c *connection) {

	if existingConn, ok := r.connections[c.deviceId]; ok && existingConn.id == c.id {
		log.Printf("Unregistering connection %s", c)
		delete(r.connections, c.deviceId)
		c.close()
		return
	}
}
