package grpc_test

import (
	"context"
	"errors"
	"fmt"
	odGrpc "github.com/nawafswe/orders-service/internal/app/orders/transport/grpc"
	"github.com/nawafswe/orders-service/internal/models"
	ordersMocks "github.com/nawafswe/orders-service/mocks/github.com/nawafswe/orders-service/pkg/v1"
	"github.com/nawafswe/orders-service/pkg/logger"
	pb "github.com/nawafswe/orders-service/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

func TestPlaceOrderService(t *testing.T) {

	tests := map[string]struct {
		Description    string
		ExpectedResult *pb.Order
		ExpectedErr    error
		Input          *pb.Order
		Assert         func(t *testing.T, useCase *ordersMocks.MockOrderUseCase, o, res *pb.Order)
	}{
		"TestSucceedPlacingOrder": {
			Description: "Should place order successfully and return a response with created order",
			ExpectedResult: &pb.Order{
				OrderId:      1,
				RestaurantId: 1,
				CustomerId:   1,
				GrandTotal:   10,
				Items: []*pb.OrderedItem{
					{
						ItemId:          1,
						OrderedItemId:   1,
						Price:           10,
						Name:            "Pepsi",
						OrderedQuantity: 1,
					},
				},
			},
			Input: &pb.Order{
				CustomerId:   1,
				RestaurantId: 1,
				GrandTotal:   10,
				Items: []*pb.OrderedItem{
					{
						OrderedItemId:   1,
						Price:           10,
						Name:            "Pepsi",
						OrderedQuantity: 1,
					},
				},
			},
		},
	}

	// follow table test pattern, define a map...
	port := 9003
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Errorf("cannot connect to server on addr: localhost:%v", port)
	}
	srv := grpc.NewServer()
	defer srv.Stop()
	orderUseCase := ordersMocks.NewMockOrderUseCase(t)
	odGrpc.NewOrderService(srv, orderUseCase, logger.NewLogger())
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Errorf("could not start a grpc server, err %v\n", err)
		}
	}()

	conn, err := grpc.Dial("localhost:9003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("failed to connect to server, err: %v \n", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Errorf("failed to close client connection, err: %v\n", err)
		}
	}(conn)
	c := pb.NewOrderServiceClient(conn)
	ctx := context.Background()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Logf("=== running %s ===", name)
			orderUseCase.On("PlaceOrder", mock.Anything, odGrpc.ToDomain(test.Input)).Return(odGrpc.ToDomain(test.Input), nil)
			res, err := c.Create(ctx, test.Input)
			// I need to handle this and return custom err
			if err == nil && test.ExpectedErr != nil {
				t.Errorf("%v failed, expected an error of %v but got %v", name, test.ExpectedErr, err)
			}
			if res == nil {
				t.Errorf("%v failed, it was expected to get a response of %v but got %v \n", name, test.ExpectedResult, res)
			}

			orderUseCase.AssertNumberOfCalls(t, "PlaceOrder", 1)
			orderUseCase.AssertCalled(t, "PlaceOrder", mock.Anything, odGrpc.ToDomain(test.Input))
			orderUseCase.AssertExpectations(t)
		})
	}
}

func TestSuccessfullyChangeOrderStatusService(t *testing.T) {
	orderUseCase := ordersMocks.NewMockOrderUseCase(t)
	port := 9009
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Errorf("failed to connect to port %d", port)
	}
	srv := grpc.NewServer()
	defer srv.Stop()
	odGrpc.NewOrderService(srv, orderUseCase, logger.NewLogger())
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Errorf("failed to start a grpc server on port %d", port)
		}
	}()
	// define client
	conn, err := grpc.Dial("localhost:9009", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error("could not establish a connection to the grpc server")
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Errorf("failed to kill client connection")
		}
	}(conn)

	c := pb.NewOrderServiceClient(conn)
	in := &pb.OrderStatus{OrderId: 1, Status: "Delivered"}
	orderUseCase.On("UpdateOrderStatus", mock.Anything, in.OrderId, in.Status).Return(models.Order{}, nil)
	_, err = c.ChangeOrderStatus(context.Background(), in)
	if err != nil {
		t.Errorf("status update field with err: %v", err)
	}
	orderUseCase.AssertNumberOfCalls(t, "UpdateOrderStatus", 1)

}

func TestFailChangeOrderStatusServiceDueInvalidOrderId(t *testing.T) {
	orderUseCase := ordersMocks.NewMockOrderUseCase(t)
	port := 9004
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Errorf("failed to connect to port %d", port)
	}
	srv := grpc.NewServer()
	defer srv.Stop()
	odGrpc.NewOrderService(srv, orderUseCase, logger.NewLogger())
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Errorf("failed to start a grpc server on port %d", port)
		}
	}()
	// define client
	conn, err := grpc.Dial("localhost:9004", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error("could not establish a connection to the grpc server")
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Errorf("failed to kill client connection")
		}
	}(conn)

	c := pb.NewOrderServiceClient(conn)
	in := &pb.OrderStatus{OrderId: -300, Status: "Delivered"}
	orderUseCase.On("UpdateOrderStatus", mock.Anything, in.OrderId, in.Status).Return(models.Order{}, errors.New("order not found"))
	_, err = c.ChangeOrderStatus(context.Background(), in)
	if err == nil {
		t.Errorf("it should fail update order status due to invalid id is passed but it did not")
	}
	orderUseCase.AssertNumberOfCalls(t, "UpdateOrderStatus", 1)

}
