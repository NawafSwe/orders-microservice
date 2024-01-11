package messaging

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/status"
)

type PUBSUB struct {
	C *pubsub.Client
}

func New(projectId string) PUBSUB {
	// we need a longed lived context to maintain client connection, using withCancel or timeout will cause unauthorized error, because the context going to be cancelled
	c, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		log.Fatalf("failed to obtain a pubsub client for project: %v, err: %v\n", projectId, err)
	}
	return PUBSUB{C: c}
}

func (p PUBSUB) CreateSub(subId string, t *pubsub.Topic) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub, err := p.C.CreateSubscription(ctx, subId, pubsub.SubscriptionConfig{
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

func (p PUBSUB) CreateTopic(topic string) *pubsub.Topic {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t := p.C.Topic(topic)
	if val, _ := t.Exists(ctx); !val {
		_, err := p.C.CreateTopic(ctx, topic)
		if err != nil {
			log.Fatalf("failed to create topic: %v, err: %v\n", topic, err)
		}
		log.Printf("topic created successfully")
	} else {
		log.Printf("topic %v already exist\n", topic)
	}
	return t
}

func (p PUBSUB) CreateTopicWithSchema(topic string, tc pubsub.TopicConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	c, _ := pubsub.NewClient(context.Background(), os.Getenv("GOOGLE_PROJECT_ID"))

	t := c.Topic(topic)
	if val, _ := t.Exists(context.Background()); !val {
		_, err := p.C.CreateTopicWithConfig(ctx, topic, &tc)
		if err != nil {
			log.Fatalf("failed to create topic: %v, err: %v\n", topic, err)
		}
		log.Printf("topic created with schema successfully")
	} else {
		log.Printf("topic %v already exist\n", topic)
	}
}

func (p PUBSUB) RetrieveTopic(topicId string) (*pubsub.Topic, error) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	t := p.C.Topic(topicId)
	log.Println(t.String())
	b, err := t.Exists(ctx)
	log.Println(b)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve topic: %v , err: %v", topicId, err)
	} else if !b {
		return nil, fmt.Errorf("topic %v is not found in your project, please create it if you are going to maintain it", topicId)
	}

	return t, nil
}
func (p PUBSUB) GetTopic(ctx context.Context, topicId string) (*pubsub.Topic, error) {
	t := p.C.Topic(topicId)

	if b, err := t.Exists(ctx); err != nil {
		return nil, fmt.Errorf("failed to get topic: %v, err: %w", topicId, err)
	} else if !b {
		return nil, fmt.Errorf("failed to get topic: %v, it is not found", topicId)
	}
	return t, nil
}
