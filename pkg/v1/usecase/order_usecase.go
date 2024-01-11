package usecase

import "C"
import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/nawafswe/orders-service/internal/models"
	"github.com/nawafswe/orders-service/pkg/messaging"
	interfaces "github.com/nawafswe/orders-service/pkg/v1"
	"github.com/nawafswe/orders-service/pkg/v1/handler/grpc"
	pb "github.com/nawafswe/orders-service/proto"
	"google.golang.org/protobuf/proto"
	"log"
)

type OrderUseCaseImpl struct {
	repo interfaces.OrderRepo
	// define an interface for messaging once you segregate business logic for publishing events and handling events from there
	pubSubClient messaging.PUBSUB
}

func New(repo interfaces.OrderRepo, ps messaging.PUBSUB) interfaces.OrderUseCase {
	return OrderUseCaseImpl{repo: repo, pubSubClient: ps}
}

func (ou OrderUseCaseImpl) PlaceOrder(ctx context.Context, order models.Order) (models.Order, error) {
	o, err := ou.repo.Create(ctx, order)
	if err != nil {
		return order, err
	}
	ou.PublishOrderCreatedEvent(ctx, grpc.FromDomain(order))
	return o, nil
}

func (ou OrderUseCaseImpl) ApproveOrder(ctx context.Context, orderId int64) error {
	return nil
}
func (ou OrderUseCaseImpl) RejectOrder(ctx context.Context, orderId int64) error {
	return nil
}
func (ou OrderUseCaseImpl) UpdateOrderStatus(ctx context.Context, orderId int64, status string) (models.Order, error) {
	return ou.repo.UpdateOrderStatus(ctx, orderId, status)
}

// Maybe Moving this logic into saga?, probably I need to do research about it

func (ou OrderUseCaseImpl) PublishOrderCreatedEvent(ctx context.Context, order *pb.Order) {
	data, err := proto.Marshal(order)
	if err != nil {
		log.Fatalf("error occured while marshling order data, err: %v\n", err)
	}
	topicId := "orderCreated"
	t, err := ou.pubSubClient.GetTopic(ctx, topicId)
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
func (ou OrderUseCaseImpl) HandleOrderApproval(ctx context.Context) {
	subId := "approveOrder"
	approveOrder := ou.pubSubClient.C.Subscription(subId)

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
		if _, err := ou.UpdateOrderStatus(ctx, orderStatus.OrderId, orderStatus.Status); err != nil {
			log.Printf("could not handle order approval, error: %v\n", err)
			msg.Nack()
			return
		}
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("Cannot receive messages for order approval at the moment, err: %v\n", err)
	}
}

func (ou OrderUseCaseImpl) HandleOrderRejection(ctx context.Context) {

}
