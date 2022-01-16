package middleware

import (
	"context"

	"google.golang.org/grpc"
)

type baseSeverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *baseSeverStream) Context() context.Context {
	return ss.ctx
}
