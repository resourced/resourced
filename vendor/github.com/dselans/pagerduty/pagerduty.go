/*
Package pagerduty is a client to the Pager Duty Integration API.
The API is described here: http://developer.pagerduty.com/documentation/integration/events

    // Create a new "trigger" event
    event := pagerduty.NewTriggerEvent(myKey, "test API")

    // Customize the incident key. If not done, pager duty will assign one to you.
    event.IncidentKey = "My Incident Key"

    // Submit the event to pager duty's API
    response, statusCode, err := pagerduty.Submit(event)
*/
package pagerduty

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// PagerDuty's API endpoint
var Endpoint = "https://events.pagerduty.com/generic/2010-04-15/create_event.json"

// See http://developer.pagerduty.com/documentation/integration/events/trigger
type Event struct {
	ServiceKey  string                 `json:"service_key"`
	EventType   string                 `json:"event_type"`
	Description string                 `json:"description"`
	IncidentKey string                 `json:"incident_key"`
	Details     map[string]interface{} `json:"details"`
}

// See http://developer.pagerduty.com/documentation/integration/events
type Response struct {
	Status      string
	Message     string
	IncidentKey string `json:"incident_key"`
	Errors      []string
}

func newEvent(serviceKey, eventType, description string) *Event {
	return &Event{ServiceKey: serviceKey, EventType: eventType, Description: description, Details: make(map[string]interface{})}
}

// Instantiates a new "trigger" Event struct
// see: http://developer.pagerduty.com/documentation/integration/events/trigger
func NewTriggerEvent(serviceKey, description string) *Event {
	return newEvent(serviceKey, "trigger", description)
}

// Instantiates a new "acknowledge" Event struct
// see: http://developer.pagerduty.com/documentation/integration/events/acknowledge
func NewAcknowledgeEvent(serviceKey, description string) *Event {
	return newEvent(serviceKey, "acknowledge", description)
}

// Instantiates a new "acknowledge" Event struct
// http://developer.pagerduty.com/documentation/integration/events/resolve
func NewResolveEvent(serviceKey, description string) *Event {
	return newEvent(serviceKey, "resolve", description)
}

// Prepare a http.Request object that can be used by the http.Client of your choice.
func PrepareRequest(event *Event) (*http.Request, error) {
	marshalled, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", Endpoint, bytes.NewBuffer(marshalled))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// Creates a Response struct from a HTTP Response of the PagerDuty's API
func NewResponse(event *Event, resp *http.Response) (*Response, error) {
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{}
	if err := json.Unmarshal(buffer, response); err != nil {
		return nil, err
	}

	if event.IncidentKey == "" {
		event.IncidentKey = response.IncidentKey
	}

	return response, nil
}

// Equivalent to a call to PrepareRequest, http.DefaultClient.Do() and NewResponse
func Submit(event *Event) (response *Response, statusCode int, err error) {
	req, err := PrepareRequest(event)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	statusCode = resp.StatusCode
	response, err = NewResponse(event, resp)
	if err != nil {
		return
	}

	return
}
