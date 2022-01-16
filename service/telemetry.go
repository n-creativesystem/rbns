package service

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/rbns/ncsfw/tracer"
)

type telemetryFunc func(ctx context.Context, spanName string, handler handler)

type handler func(ctx context.Context)

func createSpanWithPrefix(prefix string) telemetryFunc {
	return func(ctx context.Context, spanName string, handler handler) {
		createSpan(ctx, fmt.Sprintf("%s - %s", prefix, spanName), handler)
	}
}

func createSpan(ctx context.Context, spanName string, handler handler) {
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()
	handler(ctx)
}
