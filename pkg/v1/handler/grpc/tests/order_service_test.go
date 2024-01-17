package tests

import (
	"context"
	"flag"
	"fmt"
	ordersMocks "github.com/nawafswe/orders-service/mocks/github.com/nawafswe/orders-service/pkg/v1"
	orderService "github.com/nawafswe/orders-service/pkg/v1/handler/grpc"
	pb "github.com/nawafswe/orders-service/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

// mock

func TestShouldPlaceOrderSuccessfully(t *testing.T) {
	// follow table test pattern, define a map...
	port := flag.Int("port", 9003, "server port")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		t.Errorf("cannot connect to server on addr: localhost:%v", fmt.Sprintf(":%d", *port))
	}
	srv := grpc.NewServer()
	defer func() {
		srv.Stop()
	}()
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
	items := []*pb.OrderedItem{
		{
			OrderedItemId:   1,
			Price:           10,
			Sku:             "Pepsi",
			OrderedQuantity: 1,
		},
	}
	order :=
		&pb.Order{
			CustomerId: 1,
			GrandTotal: 10,
			Items:      items,
		}

	ctx := context.Background()

	// prepare mocks
	orderModel := orderService.ToDomain(order)
	orderUseCase.On("PlaceOrder", mock.Anything, orderModel).Return(orderModel, nil)

	res, err := c.Create(ctx, order)

	if err != nil {
		t.Errorf("failed to place order, err: %v\n", err)
	}

	if res == nil {
		t.Errorf("response is nil")
	}
	orderUseCase.AssertNumberOfCalls(t, "PlaceOrder", 1)
	orderUseCase.AssertCalled(t, "PlaceOrder", mock.Anything, orderModel)
	orderUseCase.AssertExpectations(t)
	if len(res.Items) != len(orderModel.Items) {
		t.Errorf("expected number of items is %v, but got %v\n", len(orderModel.Items), len(res.Items))
	}

	if res.GrandTotal != orderModel.GrandTotal {
		t.Errorf("expcpted grand total is %v, but got %v", orderModel.GrandTotal, res.GrandTotal)
	}

}
