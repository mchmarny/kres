package queue

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/mchmarny/kres/pkg/event"

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
func (r *EventRelay) Consume(e rmq.Delivery) {

	p := e.Payload()
	log.Printf("Event Payload: %v", p)

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(p), &m)
	if err != nil {
		log.Printf("Error while parsing JSON from payload: %s", err)
		e.Reject()
		return
	}

	err = r.Sender.Send(m)

	if err != nil {
		log.Printf("Error while sending event: %v", err)
		e.Reject()
	} else {
		log.Println("Acking, event sent successfully")
		e.Ack()
	}

}
