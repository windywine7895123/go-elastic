package service

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
)

func HelloService(ctx context.Context) (string, error) {
	tracer := otel.Tracer("hello-service")
	ctx, span := tracer.Start(ctx, "HelloService.Process")
	defer span.End()

	time.Sleep(50 * time.Millisecond)

	return "hello from service", nil
}
