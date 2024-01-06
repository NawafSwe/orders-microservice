package main

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
)

func GetPubSubClient() (*pubsub.Client, error) {
	googleProjectId := os.Getenv("GOOGLE_PROJECT_ID")
	return pubsub.NewClient(context.Background(), googleProjectId)

}
