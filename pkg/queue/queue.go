package queue

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/redis.v3"
	"github.com/adjust/rmq"

	"github.com/mchmarny/kres/pkg/event"
)

const (
	redisTag = "knativeEventSource"
	defaultRedisQueueName = "events"
	defaultConsumerName = "sinkRelay"
	redisHostToken = "REDIS_HOST"
	redisPassToken = "REDIS_PASS"
	redisQueueNameToken = "REDIS_QUEUE"
)

var (
	queue rmq.Queue
)

// Stop tells que (if set) to stop consuming
func Stop() {
	if queue != nil {
		queue.StopConsuming()
	}
}

// Init initializes the redis queue
func Init() error {

	host := os.Getenv(redisHostToken)
	if host == "" {
		return fmt.Errorf("Required variable undefined: %s", redisHostToken)
	}

	pass := os.Getenv(redisPassToken)
	if pass == "" {
		return fmt.Errorf("Required variable undefined: %s", redisPassToken)
	}

	queueName := os.Getenv(redisQueueNameToken)
	if queueName == "" {
		queueName = defaultRedisQueueName
	}

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       0,
	})

	log.Printf("Connecting to %s...", host)
	pong, err := client.Ping().Result()
	if err != nil {
		return fmt.Errorf("Error on PING: %v", err)
	}

	log.Printf("Success, PING=%s", pong)

	// new connection with the secure client
	log.Printf("Connecting with tag: %s...", redisTag)
	conn := rmq.OpenConnectionWithRedisClient(redisTag, client)

	if !conn.Check() {
		return fmt.Errorf("Connection failed: %v", client)
	}

	// global queue with the configured name
	log.Printf("Opening queue: %s...", queueName)
	queue = conn.OpenQueue(queueName)

	return nil

}


// ConsumeAndRelay starts consuming from queue and relaying it using sender
func ConsumeAndRelay(sender event.Sender) {

	// prefetch limit maxNumberOfEventConsumers, poll duration 1s
	// prefetchLimit == number of consumers + 1
	// this is done to avoid idling producers in times of full queues
	log.Println("Starting to consume: ...")
	queue.StartConsuming(2, time.Second)

	// add that one consumer (index 0)
	log.Printf("Adding relay consumer: %s...", defaultConsumerName)
	queue.AddConsumer(defaultConsumerName, NewEventRelay(0, sender))

	log.Println("Entering consumer loop...")
	select {}

}

