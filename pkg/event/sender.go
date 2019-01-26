package event

import (
	"bytes"
	"fmt"
	"time"
	"net/url"
 	"io/ioutil"
	"encoding/json"
    "net/http"
	"context"

	"github.com/mchmarny/knative-redis-event-source/pkg/util"
	"github.com/knative/pkg/cloudevents"
)


// SendMessages sends v02.Event based on the provided data
func NewEvent(eventType, value string) interface{} {

	now := time.Now().UTC()
	event := &v02.Event{
		SpecVersion: "0.2",
		Type:        eventType,
		Source:      *s.SourceURL,
		ID:          utils.MakeUUID(),
		Time: 		 &now,
		ContentType: "text/plain",
		Data: 		 value,
	}

	return event

}




