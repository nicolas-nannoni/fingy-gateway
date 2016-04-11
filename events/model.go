package events

import (
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

type Event struct {
	Id            uuid.UUID `id`
	ServiceId     string    `serviceId`
	CorrelationId string    `correlationId`
	Timestamp     time.Time `timestamp`
	SendTimestamp time.Time `sendTimestamp`

	Path    string      `path`
	Payload interface{} `payload`
}

// Populate extra event fields, such as its id and a creation timestamp
func (evt *Event) PrepareForSend() {
	evt.Id = uuid.NewV1()
	evt.SendTimestamp = time.Now()
}

// Verify that the event is valid and can be sent/processed
func (evt *Event) Verify() (err error) {

	if evt.Path == "" {
		return fmt.Errorf("Events should have a path set")
	}
	if evt.ServiceId == "" {
		return fmt.Errorf("Events should have a serviceId set")
	}
	return
}
