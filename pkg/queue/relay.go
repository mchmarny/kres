package queue

import (
	"log"
	"encoding/json"
	"errors"

	"github.com/mchmarny/kres/pkg/sender"
	"github.com/mchmarny/kapi/common"

	"github.com/adjust/rmq"
)

// EventRelay represents event consumer worker
type EventRelay struct {}

// NewEventRelay creates new instance of EventConsumer
func NewEventRelay() *EventRelay {
	return &EventRelay{}
}

func toStock(s string) (stock *common.SimpleStock, err error){

	if s == "" {
		return nil, errors.New("Nil json string")
	}

	var e common.SimpleStock
	if err := json.Unmarshal([]byte(s), &e); err != nil {
        log.Printf("Error while parsing JSON from payload: %s", err)
		return nil, err
    }

	return &e, nil
}

// Consume is invoked on new queue event
func (r *EventRelay) Consume(d rmq.Delivery) {

	p := d.Payload()
	log.Printf("Event Payload: %s", p)


	stock, err := toStock(p)
	if err != nil {
        log.Printf("Error converting payload to stock: %v", err)
		d.Reject()
		return
    }

	// send the raw event
	err = sender.Send(stock)

	if err != nil {
		log.Printf("Error while sending event: %v", err)
		d.Reject()
	} else {
		log.Println("Acking, event sent successfully")
		d.Ack()
	}

}
