package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/nawafswe/orders-service/orders/proto"
	domain "github.com/nawafswe/orders-service/orders/server/domain/services"
	"github.com/nawafswe/orders-service/orders/server/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ChangeOrderStatus(ctx context.Context, in *pb.OrderStatus) (*emptypb.Empty, error) {
	log.Printf("ChangeOrderStatus was invoked with in: %v\n", in)

	err := domain.ChangeOrderStatus(ctx, models.Order{
		OrderId: in.OrderId,
		Status:  in.Status,
	}, s.DB)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("order with id: %v, not found", in.OrderId))
	}

	return &emptypb.Empty{}, nil

}
