package sender

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/knative/pkg/cloudevents"

	"github.com/mchmarny/kapi/common"
)


var (
	sinkURI string
)

// Init initializes the redis queue
func Init(uri string) error {
	if uri == "" {
		return errors.New("URI required to init sender")
	}
	log.Printf("Sender initialized for %s", uri)
	sinkURI = uri
	return nil
}


// Send makes v02.Event event using passed data
// and sends it to the the provided sinkURI
func Send(e *common.SimpleStock) error {

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

	req, err := cloudevents.Binary.NewRequest(sinkURI, e, ctx)
	if err != nil {
		log.Printf("Failed to marshal: %v", err)
		return err
	}

	log.Printf("Ce-Source: %s", req.Header.Get("Ce-Source"))
	log.Printf("Ce-Eventtype: %s", req.Header.Get("Ce-Eventtype"))

	req.Header.Set("Ce-Source", "tech.knative.demo.kapi")
	req.Header.Set("Ce-Eventtype", "tech.knative.demo.kapi.stock")

	log.Printf("Posting: %v", req)

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
