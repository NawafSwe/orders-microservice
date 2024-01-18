package main

import "C"
import (
	"context"
	"flag"
	"fmt"
	"github.com/nawafswe/orders-service/internal/db"
	"github.com/nawafswe/orders-service/pkg/messaging"
	"github.com/nawafswe/orders-service/pkg/v1/repository"
	"github.com/nawafswe/orders-service/pkg/v1/usecase"
	"log"
	"net"
	"os"
	"reflect"
	"sync"

	"github.com/joho/godotenv"
	ordersGrpcService "github.com/nawafswe/orders-service/pkg/v1/handler/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	port = flag.Int("port", 9000, "gRPC server port")
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	cred, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")

	if err != nil {
		log.Fatalf("failed to create server credentials: %v\n", err)
	}

	var srvOpts []grpc.ServerOption
	srvOpts = append(srvOpts, grpc.Creds(cred))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen on addr:%v\n", lis.Addr())
	}

	s := grpc.NewServer(srvOpts...)

	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed connecting to the db, err:%v\n", err)
	}

	// generate pub sub client
	ps := messaging.New(os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		log.Fatalf("failed to connect to pub sub, err: %v\n", err)
	}
	// on main exist make sure to prevent resources leaks and close the connection of pubsub client
	defer func(service messaging.MessageService) {
		v, ok := service.(messaging.MessageServiceImpl)
		if ok {
			err := v.C.Close()
			if err != nil {
				log.Printf("failed to close the messaging client connection, err: %v\n", err)
			}
			return
		}
		log.Printf("failed to assert the type of messaging service, expected MessageServiceImpl struct but recived %v\n", reflect.TypeOf(service))
	}(ps)

	ordersRepo := repo.NewOrderRepo(dbConn)
	orderUseCase := usecase.NewOrderUseCase(ordersRepo, ps)
	ordersGrpcService.NewOrderService(s, orderUseCase)

	log.Printf("successfully connected to pub sub client...\n")
	log.Printf("Server listening at %v", lis.Addr())

	var wg sync.WaitGroup
	wg.Add(3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		defer wg.Done()
		orderUseCase.HandleOrderApproval(ctx)
	}()
	go func() {
		defer wg.Done()
		orderUseCase.HandleOrderRejection(ctx)
	}()
	go func() {
		defer wg.Done()
		// start serving requests
		if err := s.Serve(lis); err != nil {
			log.Fatalf("error ocurred when spinning a gRPC server, err: %v\n", err)
		}
		defer cancel()
	}()

	wg.Wait()
}
