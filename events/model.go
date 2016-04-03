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

func (evt *Event) PrepareForSend() {
	evt.Id = uuid.NewV1()
	evt.SendTimestamp = time.Now()
}

func (evt *Event) Verify() (err error) {

	if evt.Path == "" {
		return fmt.Errorf("Events should have a path set")
	}
	return
}
