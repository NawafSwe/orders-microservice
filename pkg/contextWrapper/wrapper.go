package contextWrapper

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"log"
)

// put this in package called contexts/wrapper.Context(ctx) do the func

// CorrelationId
// a function wraps a given context with correlation-id, if not exist before starting the process
func CorrelationId(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	log.Printf("ContextWithCorrelationId executed, current metadata: %v, corrleation-id key content: %v \n", md, md["correlation-id"])
	var correlationId string
	if ok && len(md["correlation-id"]) > 0 {
		correlationId = md["correlation-id"][0]
	} else {
		correlationId = uuid.New().String()
	}
	return context.WithValue(ctx, "correlation-id", correlationId)
}

func WithCorrelationId(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, "correlation-id", correlationId)
}
