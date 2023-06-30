package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// GetHeaderFromContext retrieves the HTTP header from the context's gRPC metadata.
func GetHeaderFromContext(ctx context.Context, key string) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	return md.Get(key)
}
