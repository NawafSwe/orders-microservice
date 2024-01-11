package grpc

import (
	"context"
	"fmt"
	"github.com/nawafswe/orders-service/internal/models"
	interfaces "github.com/nawafswe/orders-service/pkg/v1"
	pb "github.com/nawafswe/orders-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrdersServer struct {
	UseCase interfaces.OrderUseCase
	pb.UnimplementedOrderServiceServer
}

func New(s grpc.ServiceRegistrar, u interfaces.OrderUseCase) {
	pb.RegisterOrderServiceServer(s, &OrdersServer{UseCase: u})
}

func (s *OrdersServer) Create(ctx context.Context, in *pb.Order) (*pb.Order, error) {
	fmt.Printf("OrderService was invoked with Create method, with ctx: %v, in:%v\n", ctx, in)
	//return nil, status.Error(codes.Internal, "could not complete the operation")
	o, err := s.UseCase.PlaceOrder(ctx, ToDomain(in))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to place a new order, err: %w", err)
	}
	return FromDomain(o), nil

}

func (s *OrdersServer) ChangeOrderStatus(ctx context.Context, in *pb.OrderStatus) (*emptypb.Empty, error) {
	_, err := s.UseCase.UpdateOrderStatus(ctx, in.OrderId, in.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error occurred while changing order status, err: %w", err)
	}
	return &emptypb.Empty{}, nil
}
func ToDomain(o *pb.Order) models.Order {
	var items []models.OrderedItem
	for _, i := range o.Items {
		items = append(items, models.OrderedItem{
			OrderedItemId:   i.OrderedItemId,
			OrderedQuantity: i.OrderedQuantity,
			Sku:             i.Sku,
			Price:           i.Price,
		})
	}
	return models.Order{
		CustomerId: o.CustomerId,
		Status:     o.Status,
		GrandTotal: o.GrandTotal,
		Items:      items,
	}
}

func FromDomain(o models.Order) *pb.Order {
	var items []*pb.OrderedItem
	for _, i := range o.Items {
		items = append(items, &pb.OrderedItem{
			ItemId:          int64(i.ID),
			OrderedItemId:   i.OrderedItemId,
			OrderedQuantity: i.OrderedQuantity,
			Sku:             i.Sku,
			Price:           i.Price,
		})
	}
	return &pb.Order{
		OrderId:    int64(o.ID),
		CustomerId: o.CustomerId,
		Status:     o.Status,
		GrandTotal: o.GrandTotal,
		Items:      items,
	}
}
