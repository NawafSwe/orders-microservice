package main

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
)

func pullMsgs(subId string, projectId string) {

	client, err := pubsub.NewClient(context.Background(), projectId)

	if err != nil {
		log.Fatalf("error occurred while connecting to pubsub client, err: %v\n", err)
	}

	sub := client.Subscription(subId)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var received int32
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {

		log.Printf("Received message: %v\n", string(msg.Data))
		atomic.AddInt32(&received, 1)
		// acknowledge the message
		msg.Ack()
	})

	if err != nil {
		log.Fatalf("failed receiving messages, due to err: %v\n", err)
	}

	log.Printf("finished receiving messages, number of messages received: %v\n", received)

}
