package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/status"
)

func GetPubSubClient() (*pubsub.Client, error) {
	googleProjectId := os.Getenv("GOOGLE_PROJECT_ID")
	return pubsub.NewClient(context.Background(), googleProjectId)

}

func createSub(subId string, c *pubsub.Client, t *pubsub.Topic) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		log.Fatalf("failed to connect pubsub client, err: %v\n", err)
	}
	log.Printf("going to create sub\n")
	sub, err := client.CreateSubscription(ctx, subId, pubsub.SubscriptionConfig{
		Topic: t,
	})

	if err != nil {
		if e, ok := status.FromError(err); !ok {
			log.Fatalf("failed to create subscription: %v, err: %v\n", subId, err)
		} else {
			log.Printf("rpc error, err: %v", e)
			return
		}
	}

	log.Printf("Created a subscription with exactly once delivery enabled: %v\n", sub)

}
