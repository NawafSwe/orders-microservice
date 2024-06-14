package grpc

import (
	"context"
	"errors"
	"fmt"
	interfaces "github.com/nawafswe/orders-service/internal/app/orders"
	"github.com/nawafswe/orders-service/internal/models"
	"github.com/nawafswe/orders-service/pkg/logger"
	pb "github.com/nawafswe/orders-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type OrdersServer struct {
	UseCase interfaces.OrderUseCase
	pb.UnimplementedOrderServiceServer
	l logger.Logger
}

func NewOrderService(s grpc.ServiceRegistrar, u interfaces.OrderUseCase, l logger.Logger) {
	pb.RegisterOrderServiceServer(s, &OrdersServer{UseCase: u, l: l})
}

func (o *OrdersServer) Create(ctx context.Context, in *pb.Order) (*pb.Order, error) {
	start := time.Now()
	processInfo := map[string]any{
		"process":        "CreateOrder",
		"time":           start.GoString(),
		"correlation-id": ctx.Value("correlation-id"),
		"input":          in.String(),
	}
	o.l.Info(processInfo, "Executing rpc call to create order")
	//return nil, status.Error(codes.Internal, "could not complete the operation")
	if err := validateOrderCreationRequest(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	newOrder, err := o.UseCase.PlaceOrder(ctx, ToDomain(in))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to place a new order, err: %v", err)
	}
	processInfo["createdOrderId"] = newOrder.ID
	o.l.Info(processInfo, "Finished CreateOrder rpc call")
	return FromDomain(newOrder), nil

}

func (s *OrdersServer) ChangeOrderStatus(ctx context.Context, in *pb.OrderStatus) (*emptypb.Empty, error) {
	_, err := s.UseCase.UpdateOrderStatus(ctx, in.OrderId, in.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error occurred while changing order status, err: %v", err)
	}
	return &emptypb.Empty{}, nil
}
func ToDomain(o *pb.Order) models.Order {
	var items []models.OrderedItem
	for _, i := range o.Items {
		items = append(items, models.OrderedItem{
			OrderedItemId:   i.OrderedItemId,
			OrderedQuantity: i.OrderedQuantity,
			Name:            i.Name,
			Price:           i.Price,
		})
	}
	return models.Order{
		CustomerId:   o.CustomerId,
		RestaurantId: o.RestaurantId,
		Status:       o.Status,
		GrandTotal:   o.GrandTotal,
		Items:        items,
	}
}

func FromDomain(o models.Order) *pb.Order {
	var items []*pb.OrderedItem
	for _, i := range o.Items {
		items = append(items, &pb.OrderedItem{
			ItemId:          int64(i.ID),
			OrderedItemId:   i.OrderedItemId,
			OrderedQuantity: i.OrderedQuantity,
			Name:            i.Name,
			Price:           i.Price,
		})
	}
	return &pb.Order{
		OrderId:      int64(o.ID),
		CustomerId:   o.CustomerId,
		RestaurantId: o.RestaurantId,
		Status:       o.Status,
		GrandTotal:   o.GrandTotal,
		Items:        items,
	}
}

type InvalidCreateOrderRequest struct {
	Errs []error
}

func (i InvalidCreateOrderRequest) Error() string {
	return errors.Join(i.Errs...).Error()
}

// Validating order creation request before start to process it.
func validateOrderCreationRequest(o *pb.Order) error {
	var errs []error
	if o.OrderId != 0 {
		errs = append(errs, fmt.Errorf("order id should not be suplied on order placement, provided %d exepcted 0", o.OrderId))
	}

	if len(o.Items) <= 0 {
		errs = append(errs, errors.New("cannot place an order with empty items"))
	}
	for _, i := range o.Items {
		if i.OrderedItemId <= 0 {
			errs = append(errs, fmt.Errorf("ordered item id should be valid, given %d", i.OrderedItemId))
		}
		if i.ItemId != 0 {
			errs = append(errs, errors.New("item id should not be initialized"))

		}
		if i.Name == "" {
			errs = append(errs, fmt.Errorf("the name field is required"))
		}
		if i.OrderedQuantity <= 0 {
			errs = append(errs, fmt.Errorf("the quantity for item with sku %s should be greater than zero", i.Name))
		}
	}

	if o.CustomerId <= 0 {
		errs = append(errs, errors.New("the customer id must be supplied"))
	}
	if o.RestaurantId <= 0 {
		errs = append(errs, errors.New("the restaurant id must be supplied"))
	}
	if errs != nil {
		return InvalidCreateOrderRequest{Errs: errs}
	}
	return nil

}
