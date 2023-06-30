package grpc

import (
	"context"
	"time"

	"github.com/Raj63/go-sdk/logger"
	"github.com/Raj63/go-sdk/tracer"

	grpcmdw "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Client wraps up *grpc.ClientConn.
type Client struct {
	*grpc.ClientConn
}

// ClientConfig indicates how a GRPC client should be initialised.
type ClientConfig struct {
	// Address is the TCP address to connect to.
	Address string

	// KeepAlive indicates how the GRPC client should configure the connection's keep alive.
	KeepAlive struct {
		// ClientParameters is used to set keepalive parameters on the client-side. These configure
		// how the client will actively probe to notice when a connection is broken and send pings so
		// intermediaries will be aware of the liveness of the connection. Make sure these parameters
		// are set in coordination with the keepalive policy on the server, as incompatible settings
		// can result in closing of connection.
		ClientParameters struct {
			// If true, client sends keepalive pings even with no active RPCs. If false, when there are
			// no active RPCs, Time and Timeout will be ignored and no keepalive pings will be sent.
			// By default, it is false.
			PermitWithoutStream bool

			// After a duration of this time if the client doesn't see any activity it pings the server
			// to see if the transport is still alive. If set below 10s, a minimum value of 10s will be
			// used instead. By default, it is 10 * time.Second.
			Time time.Duration

			// After having pinged for keepalive check, the client waits for a duration of Timeout and
			// if no activity is seen even after that the connection is closed. By default, it is
			// 20 * time.Second.
			Timeout time.Duration
		}
	}
}

// NewClient initialises a GRPC client.
func NewClient(cfg *ClientConfig, logger *logger.Logger, tracerProvider *tracer.Provider) (*Client, error) {
	defaultClientConfig(cfg)

	conn, err := grpc.Dial(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			PermitWithoutStream: cfg.KeepAlive.ClientParameters.PermitWithoutStream,
			Time:                cfg.KeepAlive.ClientParameters.Time,
			Timeout:             cfg.KeepAlive.ClientParameters.Timeout,
		}),
		grpc.WithUnaryInterceptor(
			grpc.UnaryClientInterceptor(
				grpcmdw.ChainUnaryClient(
					otelgrpc.UnaryClientInterceptor(
						otelgrpc.WithTracerProvider(tracerProvider.TracerProvider),
					),
					grpczap.UnaryClientInterceptor(
						logger.Desugar(),
						grpczap.WithMessageProducer(loggingInterceptor),
					),
					grpczap.PayloadUnaryClientInterceptor(
						logger.Desugar(),
						func(ctx context.Context, fullMethodName string) bool {
							return true
						},
					),
				),
			),
		),
		grpc.WithStreamInterceptor(
			grpc.StreamClientInterceptor(
				grpcmdw.ChainStreamClient(
					otelgrpc.StreamClientInterceptor(
						otelgrpc.WithTracerProvider(tracerProvider.TracerProvider),
					),
					grpczap.StreamClientInterceptor(
						logger.Desugar(),
						grpczap.WithMessageProducer(loggingInterceptor),
					),
					grpczap.PayloadStreamClientInterceptor(
						logger.Desugar(),
						func(ctx context.Context, fullMethodName string) bool {
							return true
						},
					),
				),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn,
	}, nil
}

func defaultClientConfig(c *ClientConfig) *ClientConfig {
	if c.KeepAlive.ClientParameters.Time == 0 {
		c.KeepAlive.ClientParameters.Time = 10 * time.Second
	}

	if c.KeepAlive.ClientParameters.Timeout == 0 {
		c.KeepAlive.ClientParameters.Timeout = 20 * time.Second
	}

	return c
}
