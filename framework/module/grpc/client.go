package grpc

import (
	"context"
	"google.golang.org/grpc"
)

func NewClient(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, target, opts...)
}
