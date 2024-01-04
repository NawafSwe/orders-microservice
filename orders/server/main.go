package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr = "localhost:9000"

func main() {

	cred, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")

	if err != nil {
		log.Fatalf("failed to create server credentials: %v\n", err)
	}

	srvOpts := []grpc.ServerOption{}
	srvOpts = append(srvOpts, grpc.Creds(cred))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on addr:%v\n", addr)
	}
	log.Printf("Server listening on addr: %v\n", addr)

	s := grpc.NewServer(srvOpts...)
	// register server info, including services from proto buff
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error ocurred when spinning a gRPC server, err: %v\n", err)
	}

}
