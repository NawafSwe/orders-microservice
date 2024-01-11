package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	pb "github.com/nawafswe/orders-service/proto"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"time"
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

// Move this logic into order saga, or into the use cases instead as this considered a use case on its own

func (p PUBSUB) PublishOrderCreatedEvent(ctx context.Context, order *pb.Order) {
	data, err := proto.Marshal(order)
	if err != nil {
		log.Fatalf("error occured while marshling order data, err: %v\n", err)
	}
	topicId := "orderCreated"
	t, err := p.GetTopic(ctx, topicId)
	if err != nil {
		log.Fatalf("error occurred while getting topic %v, err: %w", topicId, err)
	}

	t.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	//		t.Publish(ctx, &pubsub.Message{
	//			Data: msg,
	//		})
	//		// no need to wait for publish operation
	//		//go func(result *pubsub.PublishResult) {
	//		//	ctx, cancel := context.WithCancel(context.Background())
	//		//	defer cancel()
	//		//	id, err := result.Get(ctx)
	//		//
	//		//	if err != nil {
	//		//		log.Printf("failed to publish order created event, err: %v\n", err)
	//		//	}
	//		//	log.Printf("successfully published orderCreatedEvent, msg id: %v", id)
	//		//}(result)
	//	}

}

func (p PUBSUB) HandleOrderApproval(ctx context.Context) {
	approveOrder := p.C.Subscription("approveOrder")

	if b, err := approveOrder.Exists(ctx); err != nil {
		log.Fatalf("cannot handle order approval at the moment, err: %v\n", err)
	} else if !b {
		log.Fatalf("approveOrder subscribtion resource does not exist")
	}

	err := approveOrder.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		log.Printf("recevied order approval request with msgId: %v\n", msg.ID)
		var orderStatus pb.OrderStatus
		if err := proto.Unmarshal(msg.Data, &orderStatus); err != nil {
			log.Printf("failed to unmarshal message of order status, err: %v\n", err)
			// configure nack
			msg.Nack()
			return
		}
		// update order status
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("Cannot receive messages for order approval at the moment, err: %v\n", err)
	}
}

func (p PUBSUB) HandleOrderRejection(ctx context.Context) {

}
