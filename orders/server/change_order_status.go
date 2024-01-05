package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/nawafswe/orders-service/orders/proto"
	"github.com/nawafswe/orders-service/orders/server/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ChangeOrderStatus(ctx context.Context, in *pb.OrderStatus) (*emptypb.Empty, error) {
	log.Printf("ChangeOrderStatus was invoked with in: %v\n", in)
	var order models.Order
	tx := s.DB.WithContext(ctx).Model(&order).Where("order_id=?", in.OrderId).Updates(models.Order{Status: in.Status})

	if tx.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("order with id: %v, not found", in.OrderId))
	}

	return &emptypb.Empty{}, nil

}
