// Package grpc implements a set of utilities around gRPC + Protobuf to provide low-latency
// API-contract-driven development.
package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetHeaders Retrieves all HTTP header keys as a map.
func GetHeaders(ctx context.Context, headers ...string) (map[string][]string, error) {
	var output = make(map[string][]string, len(headers))

	for _, h := range headers {
		v, err := GetHeader(ctx, h)
		if err != nil {
			return output, fmt.Errorf("unable to retrieve HTTP header=%s err=%w", h, err)
		}
		output[h] = v
	}

	return output, nil
}

// GetHeader Retrieves a HTTP header key, if available.
func GetHeader(ctx context.Context, header string) ([]string, error) {
	ids := GetHeaderFromContext(ctx, header)
	if len(ids) == 0 {
		return []string{}, status.Errorf(codes.InvalidArgument, "no %s", header)
	}

	return ids, nil
}
