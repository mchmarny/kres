package event

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/knative/pkg/cloudevents"

	"github.com/mchmarny/kapi/common"
)

// Sender represents the generic sender interface
type Sender interface {
	Send(e *common.SimpleStock) error
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
func (s *SinkSender) Send(e *common.SimpleStock) error {

	if e == nil {
		return errors.New("Required argument: event")
	}

	log.Printf("Data: %s - %s", e.ID, e.Symbol)

	ctx := cloudevents.EventContext{
		CloudEventsVersion: cloudevents.CloudEventsVersion,
		EventType:          "tech.knative.demo.kapi.stock",
		EventID:            e.ID,
		EventTime:          e.RequestOn,
		ContentType:        "application/json",
		Source:             "tech.knative.demo.kapi",
	}


	req, err := cloudevents.Binary.NewRequest(s.uri, &e, ctx)
	if err != nil {
		log.Printf("Failed to MARSHAL: %v", err)
		return err
	}

	log.Printf("Posting stock: %v", e)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusAccepted {
		log.Printf("Response Status Code: %d", resp.StatusCode)
		return fmt.Errorf("Invalid response status code: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return err
	}

	log.Printf("Response Body: %s", string(body))

	return nil

}
