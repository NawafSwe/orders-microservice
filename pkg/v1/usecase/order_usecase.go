package usecase

import "C"
import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
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
	pubSubClient messaging.MessageServiceImpl
}

func NewOrderUseCase(repo interfaces.OrderRepo, ps messaging.MessageServiceImpl) interfaces.OrderUseCase {
	return OrderUseCaseImpl{repo: repo, pubSubClient: ps}
}

func (u OrderUseCaseImpl) PlaceOrder(ctx context.Context, order models.Order) (models.Order, error) {
	for _, i := range order.Items {
		if i.OrderedQuantity <= 0 {
			return models.Order{}, fmt.Errorf("supplied quantity for item with sku %v, should be greater than zero, received is %v", i.Sku, i.OrderedQuantity)

		}
	}
	o, err := u.repo.Create(ctx, order)
	if err != nil {
		return order, err
	}
	go func() {
		u.PublishOrderCreatedEvent(grpc.FromDomain(o))
	}()

	go func() {
		u.PublishOrderStatusChanged(order)
	}()
	return o, nil
}

func (u OrderUseCaseImpl) UpdateOrderStatus(ctx context.Context, orderId int64, status string) (models.Order, error) {
	o, err := u.repo.UpdateOrderStatus(ctx, orderId, status)
	if err != nil {
		return models.Order{}, err
	}
	// the reason to use a new context here, because this function could be used by external gprc call
	// once returning to caller, the context will be canclled, to assure we resume publishing this event
	// we used long-lived context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		u.PublishOrderStatusChanged(o)
	}()
	return o, nil
}

// Maybe Moving this logic into saga?, probably I need to do research about it

func (u OrderUseCaseImpl) PublishOrderCreatedEvent(order *pb.Order) {
	data, err := proto.Marshal(order)
	if err != nil {
		log.Fatalf("error occured while marshling order data, err: %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	topicId := "orderCreated"
	t, err := u.pubSubClient.GetTopic(ctx, topicId)
	if err != nil {
		log.Fatalf("error occurred while getting topic %v, err: %v", topicId, err)
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
func (u OrderUseCaseImpl) HandleOrderApproval(ctx context.Context) {
	subId := "approveOrder"
	approveOrder, err := u.pubSubClient.GetSubscription(ctx, subId)

	if err != nil {
		log.Printf("cannot handle order approval at the moment, err: %v\n", err)
		return
	}

	err = approveOrder.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		log.Printf("recevied order approval request with msgId: %v\n", msg.ID)
		var orderStatus pb.OrderStatus
		if err := proto.Unmarshal(msg.Data, &orderStatus); err != nil {
			log.Printf("failed to unmarshal message of order status, err: %v\n", err)
			// configure nack
			msg.Nack()
			return
		}
		// update order status
		if _, err := u.UpdateOrderStatus(ctx, orderStatus.OrderId, orderStatus.Status); err != nil {
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

func (u OrderUseCaseImpl) HandleOrderRejection(ctx context.Context) {
	subId := "rejectOrder"
	s, err := u.pubSubClient.GetSubscription(ctx, subId)
	if err != nil {
		log.Printf("failed to get subscription resource, err: %v", err)
		return
	}
	err = s.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		var order pb.OrderStatus
		if err := proto.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("failed to unmarshal order status for msgId: %v, err: %v \n", msg.ID, err)
			msg.Nack()
			return
		}
		if _, err := u.UpdateOrderStatus(ctx, order.OrderId, order.Status); err != nil {
			log.Printf("failed to update order status, err: %v\n", err)
			msg.Nack()
			return
		}
		msg.Ack()
	})

	if err != nil {
		log.Fatalf("failed to receive messages for sub: %v\n", subId)
	}

}

func (u OrderUseCaseImpl) PublishOrderStatusChanged(order models.Order) {
	orderPb := pb.OrderStatus{OrderId: int64(order.ID), Status: order.Status}
	data, err := proto.Marshal(&orderPb)
	if err != nil {
		log.Printf("failed to marshal message, err: %v\n", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t, err := u.pubSubClient.GetTopic(ctx, "orderStatusChanged")
	if err != nil {
		log.Printf("could not publish event, due to err: %v\n", err)
		return
	}

	// publish status update
	t.Publish(ctx, &pubsub.Message{
		Data: data,
	})
}
