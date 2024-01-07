package messaging

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/status"
)

type PUBSUB struct {
	Client *pubsub.Client
}

func CreatePubSubClient() (*pubsub.Client, error) {
	googleProjectId := os.Getenv("GOOGLE_PROJECT_ID")
	return pubsub.NewClient(context.Background(), googleProjectId)

}

func (p PUBSUB) createSub(subId string, c *pubsub.Client, t *pubsub.Topic) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	log.Printf("going to create sub\n")
	sub, err := p.Client.CreateSubscription(ctx, subId, pubsub.SubscriptionConfig{
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

func (p PUBSUB) CreateTopic(topic string) {

	t := p.Client.Topic(topic)
	if val, _ := t.Exists(context.Background()); !val {

		_, err := p.Client.CreateTopic(context.Background(), topic)
		if err != nil {
			log.Fatalf("failed to create topic: %v, err: %v\n", topic, err)
		}
		log.Printf("topic created successfully")
	} else {
		log.Printf("topic %v already exist\n", topic)
	}

}

func (p PUBSUB) CreateTopicWithSchema(topic string, tc pubsub.TopicConfig) {
	t := p.Client.Topic(topic)
	if val, _ := t.Exists(context.Background()); !val {
		_, err := p.Client.CreateTopicWithConfig(context.Background(), topic, &tc)
		if err != nil {
			log.Fatalf("failed to create topic: %v, err: %v\n", topic, err)
		}
		log.Printf("topic created with schema successfully")
	} else {
		log.Printf("topic %v already exist\n", topic)
	}

}
