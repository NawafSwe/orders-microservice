package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	pb "github.com/nawafswe/orders-service/orders/proto"
	domain "github.com/nawafswe/orders-service/orders/server/domain/services"
	"github.com/nawafswe/orders-service/orders/server/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ChangeOrderStatus(ctx context.Context, in *pb.OrderStatus) (*emptypb.Empty, error) {
	log.Printf("ChangeOrderStatus was invoked with in: %v\n", in)
	msg, err := proto.Marshal(in)
	if err != nil {
		log.Fatalf("failed to marshal proto message")
	}
	err = domain.ChangeOrderStatus(ctx, in.OrderId, models.Order{
		Status: in.Status,
	}, s.DB)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("order with id: %v, not found", in.OrderId))
	}

	t := s.PUBSUB.Client.Topic("order_status_update")
	if val, _ := t.Exists(ctx); !val {
		log.Println("topic does not exist, going to create one...")
		t, _ = s.PUBSUB.Client.CreateTopic(ctx, "order_status_update")
	}
	// putting a correlation_id, and deal with it as a saga and pass it down to other microservices.

	t.Publish(ctx, &pubsub.Message{
		Data:       msg,
		Attributes: map[string]string{"publisher": "orders-service", "correlation_id": uuid.New().String()},
	})

	// Get will be a blocking call, it will waits till it gets the confirmation if it was published or not
	// if _, err := pr.Get(ctx); err != nil {
	// 	log.Fatalf("failed to publish a message, err:%v\n", err)
	// }

	return &emptypb.Empty{}, nil

}
