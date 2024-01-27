package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

func ServiceInterceptors(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	// validate that request has correlation id
	ctx, err := GetCorrelationIdFromRequest(ctx)
	if err != nil {
		return nil, err
	}
	// calling the handler
	h, err := handler(ctx, req)

	grpclog.Info("Request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)
	return h, err

}

func GetCorrelationIdFromRequest(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "request has invalid metadata")
	} else if len(md["correlation-id"]) == 0 {
		return nil, status.Error(codes.InvalidArgument, "request missing correlation-id")
	}
	return context.WithValue(ctx, "correlation-id", md["correlation-id"][0]), nil
}

func WithServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(ServiceInterceptors)
}
