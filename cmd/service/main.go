package main

import (
	"flag"
	"log"

	"github.com/mchmarny/kres/pkg/queue"
	"github.com/mchmarny/kres/pkg/sender"
)

var (
	sinkName  string
)

func init() {
	flag.StringVar(&sinkName, "sink", "", "Name of sink where events will be sent")
}

func main() {

	flag.Parse()

		// setup sender
	log.Printf("Initializing sender with sink: %s", sinkName)
	err := sender.Init(sinkName)
	if err != nil {
		log.Fatalf("Error while creating sink sender: %v", err)
	}

	// setup queue
	log.Println("Initializing queue")
	if err := queue.Init(); err != nil {
		log.Fatalf("Error while initializing queue: %v", err)
	}

	defer queue.Stop()

	log.Println("Starting sourcing events...")
	queue.ConsumeAndRelay()
}
