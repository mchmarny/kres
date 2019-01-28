package queue

import (
	"fmt"
	"log"
	"encoding/json"
	"time"

	"github.com/mchmarny/kres/pkg/event"
	"github.com/mchmarny/kres/pkg/common"
	"github.com/cloudevents/sdk-go/v02"

	"github.com/adjust/rmq"
)

// EventRelay represents event consumer worker
type EventRelay struct {
	Name string
	Sender event.Sender
}


// NewEventRelay creates new instance of EventConsumer
func NewEventRelay(index int, sender event.Sender) *EventRelay {
	return &EventRelay{
		Name: fmt.Sprintf("redis-event-relay-%d", index),
		Sender: sender,
	}
}

// Consume is invoked on new queue event
func (r *EventRelay) Consume(d rmq.Delivery) {

	p := d.Payload()
	log.Printf("Event Payload: %v", p)

	var e *common.SimpleStock
	if err := json.Unmarshal([]byte(p), e); err != nil {
        log.Printf("Error while parsing JSON from payload: %s", err)
		d.Reject()
		return
    }

	// send the raw event
	err := r.Sender.Send(e)

	if err != nil {
		log.Printf("Error while sending event: %v", err)
		d.Reject()
	} else {
		log.Println("Acking, event sent successfully")
		d.Ack()
	}

}
