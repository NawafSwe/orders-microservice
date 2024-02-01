package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nawafswe/orders-service/internal/db"
	"github.com/nawafswe/orders-service/internal/logger"
	"github.com/nawafswe/orders-service/pkg/messaging"
	ordersGrpcService "github.com/nawafswe/orders-service/pkg/v1/handler/grpc"
	"github.com/nawafswe/orders-service/pkg/v1/repository"
	"github.com/nawafswe/orders-service/pkg/v1/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var (
	port = flag.Int("port", 9000, "gRPC server port")
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}
	// enable tracing
	tracer.Start(
		tracer.WithEnv(os.Getenv("DD_ENV")),
		tracer.WithService(os.Getenv("SERVICE_NAME")),
		tracer.WithServiceVersion(os.Getenv("DD_VERSION")),
		tracer.WithAgentAddr(os.Getenv("DD_ADDR")),
	)
	defer tracer.Stop()

	l := logger.NewLogger()
	var srvOpts []grpc.ServerOption
	tlsEnabled := os.Getenv("TLS_ENABLED")
	log.Printf("tlsEnabled: %v\n", tlsEnabled)
	if b, _ := strconv.ParseBool(tlsEnabled); b {
		cred, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")
		srvOpts = append(srvOpts, grpc.Creds(cred))
		if err != nil {
			log.Fatalf("failed to create server credentials: %v\n", err)
		}
	}
	//l := logger.NewLogger()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen on addr:%v\n", lis.Addr())
	}
	// register middleware to intercept incoming unary requests
	srvOpts = append(srvOpts, ordersGrpcService.WithServerUnaryInterceptor())
	s := grpc.NewServer(srvOpts...)

	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed connecting to the db, err:%v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// generate pub sub client
	ps := messaging.New(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		log.Fatalf("failed to connect to pub sub, err: %v\n", err)
	}
	// on main exist make sure to prevent resources leaks and close the connection of pubsub client
	defer func(service messaging.MessageService) {
		// if failure occurred, and you want to recover and not stop the gprc server
		//defer func(){
		//if v:= recover(); v!=nil{
		//	// gracefully handle the recover, maybe log to datadog, or continue processing,
		//}
		//}()
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
	ordersGrpcService.NewOrderService(s, orderUseCase, l)

	log.Printf("successfully connected to pub sub client...\n")
	log.Printf("Server listening at %v", lis.Addr())

	var wg sync.WaitGroup
	wg.Add(3)

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
		grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
		grpclog.SetLoggerV2(grpcLog)
		l.Info(map[string]any{
			"hostname": "localhost-1",
			"appname":  "orders-service",
		}, fmt.Sprintf("service startup at %v ", time.Now().GoString()))
		// start serving requests
		if err := s.Serve(lis); err != nil {
			log.Fatalf("error ocurred when spinning a gRPC server, err: %v\n", err)
		}
		defer func() {
			log.Printf("exiting from the rpc routine\n")
		}()
	}()

	defer cancel()
	wg.Wait()
}
