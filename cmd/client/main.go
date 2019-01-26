package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"encoding/json"

	"github.com/adjust/rmq"
	"gopkg.in/redis.v3"

	"github.com/mchmarny/kres/pkg/util"

)

const (
	redisTag = "eventSourceClient"
	numberOfDeliveries = 3
)

var (
	redisHost  string
	redisPass  string
	redisQueue string
)

// Message is the content holder
type Message struct {
    ID string
	SendOn time.Time
	Value string
}

func init() {
	flag.StringVar(&redisHost, "redis", "localhost:6379", "Redis host")
	flag.StringVar(&redisPass, "password", "", "Redis password")
	flag.StringVar(&redisQueue, "queue", "events", "Redis queue")
}

func main() {

	flag.Parse()

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPass,
		DB:       0,
	})

	log.Printf("Connecting to %s...", redisHost)
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Error on PING: %v", err)
	}

	log.Printf("Success, PING=%s", pong)

	// new connection with the secure client
	log.Println("Creating client connection...")
	conn := rmq.OpenConnectionWithRedisClient(redisTag, client)

	if !conn.Check() {
		log.Fatalf("Connection failed: %v", client)
	}

	// global queue with the configured name
	log.Printf("Opening queue: %s...", redisQueue)
	queue := conn.OpenQueue(redisQueue)

	for i := 0; i < numberOfDeliveries; i++ {

		// create message
		msg := &Message{
			ID: util.MakeUUID(),
			SendOn: time.Now(),
			Value: fmt.Sprintf("delivery[%d]", i),
		}

		// msg to json
		msgBytes, _ := json.Marshal(msg)

		// send json
		ok := queue.PublishBytes(msgBytes)
		log.Printf("Delivery %d = %v", i, ok)
	}
}
