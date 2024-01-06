package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	pb "github.com/nawafswe/orders-service/orders/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/gorm"
)

var (
	port = flag.Int("port", 9000, "gRPC server port")
)

type Server struct {
	pb.OrderServiceServer
	DB *gorm.DB
}

func main() {
	err := godotenv.Load(".env")
	cred, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")

	if err != nil {
		log.Fatalf("failed to create server credentials: %v\n", err)
	}

	srvOpts := []grpc.ServerOption{}
	srvOpts = append(srvOpts, grpc.Creds(cred))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen on addr:%v\n", lis.Addr())
	}

	s := grpc.NewServer(srvOpts...)
	// register server info, including services from proto buff
	srv := &Server{}
	pb.RegisterOrderServiceServer(s, srv)
	db, err := initDB()
	if err != nil {
		log.Fatalf("failed connecting to the db, err:%v\n", err)
	}
	srv.DB = db

	log.Printf("Server listening at %v", lis.Addr())
	// start serving requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error ocurred when spinning a gRPC server, err: %v\n", err)
	}

}
