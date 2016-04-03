package services

import (
	"fmt"
	"github.com/nicolas-nannoni/fingy-server/events"
	"github.com/parnurzeal/gorequest"
	"net/url"
)

type Service struct {
	Id   string
	Host string
	Port uint
}

type registry struct {
	services map[string]*Service
}

var Registry = registry{services: make(map[string]*Service)}

func (r *registry) Dispatch(evt *events.Event) (resp *events.Event, err error) {

	service := r.getService(evt)
	if service == nil {
		return nil, fmt.Errorf("Unknown service %s", service.Id)
	}

	body, errs := service.send(evt)
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

func (r *registry) RegisterService(service *Service) {
	r.services[service.Id] = service
}

func (r *registry) getService(evt *events.Event) (service *Service) {
	service = r.services[evt.ServiceId]
	return
}

func (s *Service) send(evt *events.Event) (response string, errs []error) {

	request := gorequest.New()
	u := url.URL{Scheme: "http", Host: s.Host, Path: evt.Path}
	_, body, errs := request.Get(u.String()).End()

	return body, errs
}
