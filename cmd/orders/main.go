package main

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
	client := messaging.New(os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		log.Fatalf("failed to connect to pub sub, err: %v\n", err)
	}
	defer client.C.Close()

	ordersRepo := repo.New(dbConn)
	orderUseCase := usecase.New(ordersRepo, client)
	ordersGrpcService.New(s, orderUseCase)

	log.Printf("successfully connected to pub sub client...\n")
	log.Printf("Server listening at %v", lis.Addr())

	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		defer wg.Done()
		orderUseCase.HandleOrderApproval(ctx)
	}()

	go func() {
		// start serving requests
		if err := s.Serve(lis); err != nil {
			log.Fatalf("error ocurred when spinning a gRPC server, err: %v\n", err)
		}
		defer cancel()
	}()

	wg.Wait()
}
