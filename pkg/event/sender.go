package event

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	ce "github.com/knative/pkg/cloudevents"
	"github.com/mchmarny/kres/pkg/util"
)

const (
	redisEventSource = "io.redis"
	redisEventType   = "io.redis.queue"
)

// Sender represents the generic sender interface
type Sender interface {
	Send(data map[string]interface{}) error
}

// SinkSender sends messages to sink
type SinkSender struct {
	uri string
}

// NewSinkSender creates configured instance of the sink sender
func NewSinkSender(sinkURI string) (sender Sender, err error) {

	if sinkURI == "" {
		return nil, errors.New("Required argument: sinkURI")
	}

	return &SinkSender{
		uri: sinkURI,
	}, nil
}

// Send makes v02.Event event using passed data
// and sends it to the the provided sinkURI
func (s *SinkSender) Send(data map[string]interface{}) error {

	if data == nil {
		return errors.New("Required argument: data")
	}

	log.Printf("Data: %v", data)

	ex := ce.EventContext{
		CloudEventsVersion: ce.CloudEventsVersion,
		EventType:          redisEventType,
		EventTypeVersion:   "v0.1",
		EventID:            util.MakeUUID(),
		EventTime:          time.Now(),
		ContentType:        "application/json",
		Source:             redisEventSource,
	}

	req, err := ce.Binary.NewRequest(s.uri, data, ex)
	if err != nil {
		log.Printf("Error creating new quest: %v", err)
		return err
	}

	log.Printf("Posting to %s: %v", s.uri, data)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusAccepted {
		log.Printf("Response Status: %s", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Response Body: %s", string(body))
		return fmt.Errorf("Invalid response status code: %s", resp.Status)
	}

	return nil

}
