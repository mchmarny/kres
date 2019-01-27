package event

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cloudevents/sdk-go/v02"
)

const (
	redisEventSource = "io.redis"
	redisEventType   = "io.redis.queue"
)

// Sender represents the generic sender interface
type Sender interface {
	Send(e *v02.Event) error
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
func (s *SinkSender) Send(e *v02.Event) error {

	if e == nil {
		return errors.New("Required argument: event")
	}

	log.Printf("Data: %v", e)

	m := v02.NewDefaultHTTPMarshaller()
	var req *http.Request
	err := m.ToRequest(req, e)
	if err != nil {
		log.Printf("Unable to marshal event into http Request: %v", err)
		return err
	}

	log.Printf("Posting to %s: %v", s.uri, req)

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
