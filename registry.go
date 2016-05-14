package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/nicolas-nannoni/fingy-gateway/events"
	"github.com/parnurzeal/gorequest"
	"net/url"
)

type Service struct {
	Id   string
	Host string
	Port uint

	deviceRegistry map[string]*connection
}

type registry struct {
	services map[string]*Service
}

var Registry = registry{
	services: make(map[string]*Service),
}

// Dispatch an Event coming from a given deviceId to the appropriate serviceId that exists in the Registry
func (r *registry) Dispatch(serviceId string, deviceId string, evt *events.Event) (resp *events.Event, err error) {

	service := r.services[serviceId]
	if service == nil {
		return nil, fmt.Errorf("Unknown service %s", evt.ServiceId)
	}

	body, errs := service.send(evt, deviceId)
	if errs != nil {
		return nil, fmt.Errorf("Error while contacting service %s: %v", service.Id, errs)
	}

	resp = &events.Event{
		ServiceId:     service.Id,
		CorrelationId: evt.Id.String(),
		Path:          evt.Path,
		Payload:       body,
	}

	return resp, nil
}

// Send an event to the given deviceId (registered in the given service)
func (r *registry) SendToDevice(serviceId string, deviceId string, evt *events.Event) (err error) {

	s := r.services[serviceId]
	if s == nil {
		err = fmt.Errorf("Unknown service %s. Message sending to %s aborted", serviceId, deviceId)
		return
	}
	return s.SendToDevice(deviceId, evt)
}

// Send an event to the given deviceId (registered in the current service)
func (s *Service) SendToDevice(deviceId string, evt *events.Event) (err error) {

	c, ok := s.deviceRegistry[deviceId]
	if !ok {
		return fmt.Errorf("No connection to device %s for service %s", deviceId, s)
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
	log.Debugf("Pushing message %s to send queue of %s", msg, c)
	c.send <- msg

	return
}

// Add a service to the Fingy Service registry
func (r *registry) RegisterService(service *Service) {
	r.services[service.Id] = service
	service.deviceRegistry = make(map[string]*connection)

}

// Lookup the Service entity matching the given event
func (r *registry) getServiceForEvent(evt *events.Event) (service *Service) {
	service = r.services[evt.ServiceId]
	return
}

// Lookup the Service entity matching the given connection
func (r *registry) getServiceForConnection(c *connection) (service *Service) {
	service = r.services[c.serviceId]
	return
}

// Send the HTTP request towards the final service
func (s *Service) send(evt *events.Event, deviceId string) (response string, errs []error) {

	request := gorequest.New()
	u := url.URL{Scheme: "http", Host: s.Host, Path: fmt.Sprintf("/devices/%s%s", deviceId, evt.Path)}
	resp, body, errs := request.Get(u.String()).End()

	log.Debugf("Response from service %s: %s", s, resp)

	return body, errs
}

// Register a connection (device)
func (r *registry) registerConnection(c *connection) (err error) {

	log.Infof("Registering connection %s", c)

	s := r.getServiceForConnection(c)
	if s == nil {
		err = fmt.Errorf("Unknown service with id %s, unable to register connection", c.serviceId)
		return
	}

	if existingConn, ok := s.deviceRegistry[c.deviceId]; ok {
		log.Debugf("Existing registration for device %s. Closing old connection %s", c.deviceId, existingConn)
		r.unregisterConnection(existingConn)
	}

	s.deviceRegistry[c.deviceId] = c
	return
}

// Unregister a connection (device)
func (r *registry) unregisterConnection(c *connection) (err error) {

	s := r.getServiceForConnection(c)
	if s == nil {
		err = fmt.Errorf("Unknown service with id %s, unable to unregister connection", c.serviceId)
		return
	}

	if existingConn, ok := s.deviceRegistry[c.deviceId]; ok && existingConn.id == c.id {
		log.Infof("Unregistering connection %s", c)
		delete(s.deviceRegistry, c.deviceId)
		c.Close()
		return
	}

	return
}
