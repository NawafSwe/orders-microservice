package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/nawafswe/orders-service/orders/server/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/nawafswe/orders-service/orders/proto"
)

func (s *Server) Create(ctx context.Context, in *pb.Order) (*pb.Order, error) {
	fmt.Printf("OrderService was invoked with Create method, with ctx: %v, in:%v\n", ctx, in)
	createdOrder := models.Order{
		CustomerId: in.CustomerId,
		GrandTotal: in.GrandTotal,
		Status:     "new",
	}
	tx := s.DB.WithContext(ctx).Create(&createdOrder)

	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit order, err: %v", tx.Error)
	}
	if tx.RowsAffected > 0 {
		log.Printf("Created order ===> %v\n", createdOrder)
		log.Printf("order id  ===> %v\n", createdOrder.ID)
		// create items
		var createdItems []models.OrderedItem
		for _, item := range in.Items {
			createdItems = append(createdItems, models.OrderedItem{
				OrderedQuantity: item.OrderedQuantity,
				Price:           item.Price,
				Sku:             item.Sku,
				OrderID:         createdOrder.ID,
			})
		}
		log.Printf("createdItems ===> %v \n", createdItems)
		btx := s.DB.WithContext(ctx).CreateInBatches(&createdItems, 10)

		if err := btx.Error; err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "failed to create items of order, err: %v", err)
		}
		var preparedItems []*pb.OrderedItem
		for _, item := range createdItems {
			preparedItems = append(preparedItems, &pb.OrderedItem{
				ItemId:          int64(item.ID),
				OrderedQuantity: item.OrderedQuantity,
				Price:           item.Price,
				Sku:             item.Sku,
			})
		}
		pbOrder := &pb.Order{
			OrderId:    int64(createdOrder.ID),
			Items:      preparedItems,
			GrandTotal: createdOrder.GrandTotal,
			Status:     createdOrder.Status,
			CustomerId: createdOrder.CustomerId,
		}
		publishOrderCreatedEvent(ctx, s, pbOrder.ProtoReflect().Interface())
		return pbOrder, nil
	}
	return nil, status.Error(codes.Internal, "could not complete the operation")

}

func publishOrderCreatedEvent(ctx context.Context, s *Server, data proto.Message) {
	msg, err := proto.Marshal(data)
	if err != nil {
		log.Fatalf("failed to marshal proto message")
	}
	orderCreatedTopic := "orderCreated"
	t := s.PUBSUB.Client.Topic(orderCreatedTopic)
	if b, err := t.Exists(ctx); err != nil {
		log.Printf("failed to publish for topic: %v, err: %v\n", orderCreatedTopic, err)
	} else if !b {
		log.Printf("failed to publish for topic %v, due to topic does not exist\n", orderCreatedTopic)
	} else {
		t.Publish(ctx, &pubsub.Message{
			Data: msg,
		})

		//go func(result *pubsub.PublishResult) {
		//	ctx, cancel := context.WithCancel(context.Background())
		//	defer cancel()
		//	id, err := result.Get(ctx)
		//
		//	if err != nil {
		//		log.Printf("failed to publish order created event, err: %v\n", err)
		//	}
		//	log.Printf("successfully published orderCreatedEvent, msg id: %v", id)
		//}(result)
	}
}
