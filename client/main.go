package main

import (
	pb "github.com/nawafswe/orders-service/proto"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load .env file")
	}

	tlsEnabled := os.Getenv("TLS_ENABLED")
	log.Printf("tlsEnabled: %v\n", tlsEnabled)
	var copts []grpc.DialOption
	if b, _ := strconv.ParseBool(tlsEnabled); b {
		cred, err := credentials.NewClientTLSFromFile("ssl/server.crt", "")
		copts = append(copts, grpc.WithTransportCredentials(cred))
		if err != nil {
			log.Fatalf("failed to obtain credentials, err: %v\n", err)
		}
	} else {
		copts = append(copts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
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
