package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nawafswe/orders-service/orders/server/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		createdItems := []models.OrderedItem{}
		for _, item := range in.Items {
			createdItems = append(createdItems, models.OrderedItem{
				OrderedQuantity: item.OrderedQuantity,
				Price:           item.Price,
				Sku:             item.Sku,
				OrderID:         uint(createdOrder.ID),
			})
		}
		log.Printf("createdItems ===> %v \n", createdItems)
		btx := s.DB.WithContext(ctx).CreateInBatches(&createdItems, 10)

		if err := btx.Error; err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "failed to create items of order, err: %v", err)
		}
		preparedItems := []*pb.OrderedItem{}
		for _, item := range createdItems {
			preparedItems = append(preparedItems, &pb.OrderedItem{
				ItemId:          int64(item.ID),
				OrderedQuantity: item.OrderedQuantity,
				Price:           item.Price,
				Sku:             item.Sku,
			})
		}
		return &pb.Order{
			OrderId:    int64(createdOrder.ID),
			Items:      preparedItems,
			GrandTotal: createdOrder.GrandTotal,
			Status:     createdOrder.Status,
			CustomerId: createdOrder.CustomerId,
		}, nil
	}
	return nil, status.Error(codes.Internal, "could not complete the operation")

}
