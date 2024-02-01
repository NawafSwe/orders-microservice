package usecase

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/nawafswe/orders-service/contextWrapper"
	"github.com/nawafswe/orders-service/internal/models"
	"github.com/nawafswe/orders-service/pkg/logger"
	"github.com/nawafswe/orders-service/pkg/messaging"
	interfaces "github.com/nawafswe/orders-service/pkg/v1"
	ordersService "github.com/nawafswe/orders-service/pkg/v1/handler/grpc"
	pb "github.com/nawafswe/orders-service/proto"
	"google.golang.org/protobuf/proto"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"log"
	"os"
	"time"
)

type OrderUseCaseImpl struct {
	repo interfaces.OrderRepo
	// define an interface for messaging once you segregate business logic for publishing events and handling events from there
	pubSubClient messaging.MessageService
	l            logger.Logger
}

func NewOrderUseCase(repo interfaces.OrderRepo, ps messaging.MessageService, l logger.Logger) interfaces.OrderUseCase {
	return OrderUseCaseImpl{repo: repo, pubSubClient: ps, l: l}
}

func (u OrderUseCaseImpl) PlaceOrder(ctx context.Context, order models.Order) (models.Order, error) {
	for _, i := range order.Items {
		if i.OrderedQuantity <= 0 {
			return models.Order{}, fmt.Errorf("supplied quantity for item with name %v, should be greater than zero, received is %v", i.Name, i.OrderedQuantity)

		}
	}
	o, err := u.repo.Create(ctx, order)
	if err != nil {
		return models.Order{}, err
	}
	ctx = contextWrapper.CorrelationId(ctx)
	u.PublishOrderCreatedEvent(ctx, o)
	u.PublishOrderStatusChanged(ctx, order)
	return o, nil
}

func (u OrderUseCaseImpl) UpdateOrderStatus(ctx context.Context, orderId int64, status string) (models.Order, error) {
	if status == "" {
		return models.Order{}, models.InvalidStatusChangeErr{Message: fmt.Sprintf("given status '%v' is invalid", status)}
	}
	o, err := u.repo.UpdateOrderStatus(ctx, orderId, status)
	if err != nil {
		return models.Order{}, err
	}
	u.PublishOrderStatusChanged(ctx, o)

	return o, nil
}

// Maybe Moving this logic into saga?, probably I need to do research about it

func (u OrderUseCaseImpl) PublishOrderCreatedEvent(ctx context.Context, order models.Order) {
	data, err := proto.Marshal(ordersService.FromDomain(order))
	if err != nil {
		log.Printf("error occured while marshling order data, err: %v\n", err)
	}
	msgId, ok := ctx.Value("correlation-id").(string)
	if !ok {
		ctx = contextWrapper.CorrelationId(ctx)
		msgId = ctx.Value("correlation-id").(string)
	}

	u.pubSubClient.PublishAsync(ctx, "orderCreated", &pubsub.Message{Data: data, Attributes: map[string]string{"correlation-id": msgId, "service": "orders", "host": "localhost"}})

}
func (u OrderUseCaseImpl) HandleOrderApproval(ctx context.Context) {
	subId := "approveOrder"
	processName := "HandleOrderApproval"
	approveOrder, err := u.pubSubClient.GetSubscription(ctx, subId)

	if err != nil {
		log.Printf("cannot handle order approval at the moment, err: %v\n", err)
		return
	}
	u.l.Info(map[string]any{
		"process": processName,
		"subId":   subId,
		"time":    time.Now(),
	}, "starting to handling order approvals requests")

	err = approveOrder.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		msgId := msg.Attributes["correlation-id"]
		start := time.Now()
		span, _ := tracer.StartSpanFromContext(ctx, processName)
		defer span.Finish()

		var orderStatus pb.OrderStatus
		if err := proto.Unmarshal(msg.Data, &orderStatus); err != nil {
			log.Printf("failed to unmarshal message of order status, err: %v\n", err)
			// configure nack
			msg.Nack()
			return
		}
		if correlationId, ok := msg.Attributes["correlation-id"]; ok {
			ctx = contextWrapper.WithCorrelationId(ctx, correlationId)
		}
		u.l.Info(map[string]any{
			"process":        processName,
			"context":        fmt.Sprintf("Handling HandleOrderApproval for orderId: %v, corrleation-id: %v", orderStatus.OrderId, msgId),
			"correlation-id": msgId,
			"service":        os.Getenv("SERVICE_NAME"),
			"time":           time.Since(start),
		}, "Processing incoming order")
		// update order status
		processedOrder, err := u.UpdateOrderStatus(ctx, orderStatus.OrderId, orderStatus.Status)
		if err != nil {
			log.Printf("could not handle order approval for order: %v, error: %v\n", orderStatus.OrderId, err)
			msg.Nack()
			return
		}
		u.l.Info(map[string]any{
			"process":        processName,
			"context":        fmt.Sprintf("restaurantId: %v approved the incoming order %v succesfully", processedOrder.RestaurantId, processedOrder.ID),
			"correlation-id": msgId,
			"service":        os.Getenv("SERVICE_NAME"),
			"time":           time.Since(start),
		}, "Order processed")
		msg.Ack()
	})
	if err != nil {
		log.Printf("Cannot receive messages for order approval at the moment, err: %v\n", err)
	}
}

func (u OrderUseCaseImpl) HandleOrderRejection(ctx context.Context) {
	subId := "rejectOrder"
	s, err := u.pubSubClient.GetSubscription(ctx, subId)
	if err != nil {
		log.Printf("failed to get subscription resource, err: %v", err)
		return
	}
	log.Println("===== starting to handle order rejection events =====")
	err = s.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		var order pb.OrderStatus
		if err := proto.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("failed to unmarshal order status for msgId: %v, err: %v \n", msg.ID, err)
			msg.Nack()
			return
		}
		if correlationId, ok := msg.Attributes["correlation-id"]; ok {
			ctx = contextWrapper.WithCorrelationId(ctx, correlationId)
		}
		if _, err := u.UpdateOrderStatus(ctx, order.OrderId, order.Status); err != nil {
			log.Printf("failed to update order status, err: %v\n", err)
			msg.Nack()
			return
		}
		msg.Ack()
	})

	if err != nil {
		log.Printf("failed to receive messages for sub: %v\n", subId)
	}

}

func (u OrderUseCaseImpl) PublishOrderStatusChanged(ctx context.Context, order models.Order) {
	orderPb := pb.OrderStatus{OrderId: int64(order.ID), Status: order.Status}
	data, err := proto.Marshal(&orderPb)
	if err != nil {
		log.Printf("failed to marshal message, err: %v\n", err)
		return
	}
	msgId, ok := ctx.Value("correlation-id").(string)
	if !ok {
		ctx = contextWrapper.CorrelationId(ctx)
		msgId = ctx.Value("correlation-id").(string)
	}
	u.pubSubClient.PublishAsync(ctx, "orderStatusChanged", &pubsub.Message{
		Data:       data,
		Attributes: map[string]string{"correlation-id": msgId, "service": "orders", "host": "localhost"},
	})

}
