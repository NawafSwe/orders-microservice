package main

import (
	pb "github.com/nawafswe/orders-service/proto"
	"log"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load .env file")
	}

	cred, err := credentials.NewClientTLSFromFile("ssl/server.crt", "")

	if err != nil {
		log.Fatalf("failed to obtain credentials, err: %v\n", err)
	}
	copts := []grpc.DialOption{}
	copts = append(copts, grpc.WithTransportCredentials(cred))
	addr := "localhost:9000"
	conn, err := grpc.Dial(addr, copts...)
	if err != nil {
		log.Fatalf("failed to connect to addr: %v, err: %v \n", addr, err)
	}

	defer conn.Close()

	c := pb.NewOrderServiceClient(conn)
	createOrder(c)
	// changeOrderStatus(c)
}
