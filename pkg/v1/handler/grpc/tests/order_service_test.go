package grpc_tests

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/nawafswe/orders-service/internal/models"
	ordersMocks "github.com/nawafswe/orders-service/mocks/github.com/nawafswe/orders-service/pkg/v1"
	orderService "github.com/nawafswe/orders-service/pkg/v1/handler/grpc"
	pb "github.com/nawafswe/orders-service/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net"
	"testing"
)

func TestPlaceOrderService(t *testing.T) {

	tests := map[string]struct {
		Description    string
		ExpectedResult func(o *pb.Order) (*pb.Order, error)
		Input          *pb.Order
		Before         func(useCase *ordersMocks.MockOrderUseCase, order *pb.Order)
		Assert         func(t *testing.T, useCase *ordersMocks.MockOrderUseCase, o, res *pb.Order)
	}{
		"TestSucceedPlacingOrder": {
			Description: "Should place order successfully and return a response with created order",
			ExpectedResult: func(o *pb.Order) (*pb.Order, error) {
				orderModel := orderService.ToDomain(o)
				return orderService.FromDomain(orderModel), nil
			},
			Input: &pb.Order{
				CustomerId: 1,
				GrandTotal: 10,
				Items: []*pb.OrderedItem{
					{
						OrderedItemId:   1,
						Price:           10,
						Sku:             "Pepsi",
						OrderedQuantity: 1,
					},
				},
			},
			Before: func(useCase *ordersMocks.MockOrderUseCase, o *pb.Order) {
				useCase.On("PlaceOrder", mock.Anything, orderService.ToDomain(o)).Return(orderService.ToDomain(o), nil)
			},
			Assert: func(t *testing.T, useCase *ordersMocks.MockOrderUseCase, in, res *pb.Order) {
				useCase.AssertNumberOfCalls(t, "PlaceOrder", 1)
				useCase.AssertCalled(t, "PlaceOrder", mock.Anything, orderService.ToDomain(in))
				useCase.AssertExpectations(t)
				orderModel := orderService.ToDomain(in)
				if len(res.Items) != len(orderModel.Items) {
					t.Errorf("expected number of items is %v, but got %v\n", len(orderModel.Items), len(res.Items))
				}

				if res.GrandTotal != orderModel.GrandTotal {
					t.Errorf("expcpted grand total is %v, but got %v", orderModel.GrandTotal, res.GrandTotal)
				}
			},
		},
		"TestFailPlaceOrder": {
			Description: "Should fail place order and return a nil response with error, due to invalid quantity passed",
			ExpectedResult: func(o *pb.Order) (*pb.Order, error) {
				orderModel := orderService.ToDomain(o)
				return orderService.FromDomain(orderModel), nil
			},
			Input: &pb.Order{
				CustomerId: 1,
				GrandTotal: 10,
				Items: []*pb.OrderedItem{
					{
						OrderedItemId:   1,
						Price:           10,
						Sku:             "Pepsi",
						OrderedQuantity: -1,
					},
				},
			},
			Before: func(useCase *ordersMocks.MockOrderUseCase, o *pb.Order) {
				err := status.Errorf(codes.Internal, "supplied quantity for item with sku %v, should be greater than zero, received is %v", o.Items[0].Sku, o.Items[0].OrderedQuantity)
				useCase.On("PlaceOrder", mock.Anything, orderService.ToDomain(o)).Return(orderService.ToDomain(&pb.Order{}), err)
			},
			Assert: func(t *testing.T, useCase *ordersMocks.MockOrderUseCase, in, res *pb.Order) {
				useCase.AssertNumberOfCalls(t, "PlaceOrder", 1)
				useCase.AssertCalled(t, "PlaceOrder", mock.Anything, orderService.ToDomain(in))
				useCase.AssertExpectations(t)

			},
		},
	}

	// follow table test pattern, define a map...
	port := flag.Int("port", 9003, "server port")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		t.Errorf("cannot connect to server on addr: localhost:%v", fmt.Sprintf(":%d", *port))
	}
	srv := grpc.NewServer()
	defer srv.Stop()
	orderUseCase := ordersMocks.NewMockOrderUseCase(t)
	orderService.NewOrderService(srv, orderUseCase)
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
		// skipping, because test driven table is not suitable here, where each test case should have
		// isolated mocks with it.
		t.Skipf("TestFailPlaceOrder")
		t.Run(name, func(t *testing.T) {
			t.Logf("=== running %s ===", name)
			test.Before(orderUseCase, test.Input)
			res, err := c.Create(ctx, test.Input)
			expectedRes, expectedErr := test.ExpectedResult(test.Input)
			// I need to handle this and return custom err
			if errors.Is(err, expectedErr) {
				t.Errorf("%v failed, expected an error of %v but got %v", name, expectedErr, err)
			}
			if res == nil {
				t.Errorf("%v failed, it was expected to get a response of %v but got %v \n", name, expectedRes, res)
			}

			test.Assert(t, orderUseCase, test.Input, res)
		})
	}
}

func TestSuccessfullyChangeOrderStatusService(t *testing.T) {
	orderUseCase := ordersMocks.NewMockOrderUseCase(t)
	port := flag.Int("port", 9003, "server port")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		t.Errorf("failed to connect to port %d", *port)
	}
	srv := grpc.NewServer()
	defer srv.Stop()
	orderService.NewOrderService(srv, orderUseCase)
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Errorf("failed to start a grpc server on port %d", *port)
		}
	}()
	// define client
	conn, err := grpc.Dial("localhost:9003", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	port := flag.Int("port", 9003, "server port")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		t.Errorf("failed to connect to port %d", *port)
	}
	srv := grpc.NewServer()
	defer srv.Stop()
	orderService.NewOrderService(srv, orderUseCase)
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Errorf("failed to start a grpc server on port %d", *port)
		}
	}()
	// define client
	conn, err := grpc.Dial("localhost:9003", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
