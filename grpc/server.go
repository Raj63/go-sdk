package grpc

import (
	"context"
	"net"
	"time"

	"github.com/Raj63/go-sdk/logger"
	"github.com/Raj63/go-sdk/tracer"

	grpcmdw "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server wraps up *grpc.Server.
type Server struct {
	config           *ServerConfig
	listener         net.Listener
	logger           *logger.Logger
	server           *grpc.Server
	preStartCallback func() error
}

// ServerConfig indicates how a gRPC server should be initialised.
type ServerConfig struct {
	// Address is the TCP address to listen on.
	Address string

	// GracefulShutdownHandler is a function that runs before the gRPC server is gracefully shut down.
	GracefulShutdownHandler func() error

	// GRPCGatewayServer indicates the `grpc-gateway` server that the service is connected to.
	GRPCGatewayServer string

	// KeepAlive indicates how the gRPC server should configure the connection's keep alive.
	KeepAlive struct {
		// EnforcementPolicy is used to set keepalive enforcement policy on the server-side. Server
		// will close connection with a client that violates this policy.
		EnforcementPolicy struct {
			// MinTime is the minimum amount of time a client should wait before sending a keepalive
			// ping. By default, it is 5 * time.Second.
			MinTime time.Duration

			// If true, server allows keepalive pings even when there are no active streams(RPCs). If
			// false, and client sends ping when there are no active streams, server will send GOAWAY
			// and close the connection. By default, it is false.
			PermitWithoutStream bool
		}

		// ServerParameters is used to set keepalive and max-age parameters on the server-side.
		ServerParameters struct {
			// MaxConnectionAge is a duration for the maximum amount of time a connection may exist
			// before it will be closed by sending a GoAway. A random jitter of +/-10% will be added
			// to MaxConnectionAge to spread out connection storms.
			MaxConnectionAge time.Duration

			// MaxConnectionIdle is a duration for the amount of time after which an idle connection
			// would be closed by sending a GoAway. Idleness duration is defined since the most recent
			// time the number of outstanding RPCs became zero or the connection establishment.
			MaxConnectionIdle time.Duration

			// After a duration of this time if the server doesn't see any activity it pings the client
			// to see if the transport is still alive. If set below 1s, a minimum value of 1s will be
			// used instead.
			Time time.Duration

			// After having pinged for keepalive check, the server waits for a duration of Timeout and
			// if no activity is seen even after that the connection is closed.
			Timeout time.Duration
		}
	}

	// ReflectionService indicates if the gRPC server should register the reflection service.
	ReflectionService bool

	// TracerProvider is the provider that uses the exporter to push traces to the collector.
	TracerProvider *tracer.Provider

	// TracerProviderShutdownHandler is a function that shuts down the tracer's exporter/provider before
	// the gRPC server is gracefully shut down.
	TracerProviderShutdownHandler func() error
}

var (
	system = "ok" // ok string represents the health of the system
)

// NewServer initialises a gRPC server.
func NewServer(c *ServerConfig, logger *logger.Logger, preStartCallback func() error, interceptors ...grpc.UnaryServerInterceptor) (*Server, error) {
	defaultServerConfig(c)
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		grpcrecovery.UnaryServerInterceptor(),
		otelgrpc.UnaryServerInterceptor(
			otelgrpc.WithTracerProvider(c.TracerProvider),
		),
		grpczap.UnaryServerInterceptor(
			logger.Desugar(),
			grpczap.WithMessageProducer(loggingInterceptor),
		),
		grpczap.PayloadUnaryServerInterceptor(
			logger.Desugar(),
			func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
				return true
			},
		),
	}
	interceptors = append(defaultInterceptors, interceptors...)
	srv := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             c.KeepAlive.EnforcementPolicy.MinTime,
				PermitWithoutStream: c.KeepAlive.EnforcementPolicy.PermitWithoutStream,
			},
		),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionAge:  c.KeepAlive.ServerParameters.MaxConnectionAge,
				MaxConnectionIdle: c.KeepAlive.ServerParameters.MaxConnectionIdle,
				Time:              c.KeepAlive.ServerParameters.Time,
				Timeout:           c.KeepAlive.ServerParameters.Timeout,
			},
		),
		grpc.UnaryInterceptor(
			grpcmdw.ChainUnaryServer(interceptors...),
		),
		grpc.StreamInterceptor(
			grpcmdw.ChainStreamServer(
				grpcrecovery.StreamServerInterceptor(),
				otelgrpc.StreamServerInterceptor(
					otelgrpc.WithTracerProvider(c.TracerProvider),
				),
				grpczap.StreamServerInterceptor(
					logger.Desugar(),
					grpczap.WithMessageProducer(loggingInterceptor),
				),
				grpczap.PayloadStreamServerInterceptor(
					logger.Desugar(),
					func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
						return true
					},
				),
			),
		),
	)

	healthcheck := health.NewServer()
	healthcheck.SetServingStatus(system, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(srv, healthcheck)

	if c.ReflectionService {
		reflection.Register(srv)
	}

	return &Server{
		c,
		nil,
		logger,
		srv,
		preStartCallback,
	}, nil
}

// Addr returns the server's network address.
func (s *Server) Addr() string {
	return s.config.Address
}

// GRPCServer returns the internal gRPC server instance.
func (s *Server) GRPCServer() *grpc.Server {
	return s.server
}

// GracefulShutdownHandler is a function that runs before the gRPC server is gracefully shut down.
func (s *Server) GracefulShutdownHandler() error {
	return s.config.GracefulShutdownHandler()
}

// GracefulStop stops the gRPC server gracefully. It stops the server from accepting new
// connections and RPCs and blocks until all the pending RPCs are finished.
func (s *Server) GracefulStop() error {
	s.server.GracefulStop()

	return nil
}

// GRPCGatewayClients indicates the services that the server is proxying via `grpc-gateway`. As a
// gRPC server isn't meant to proxy any HTTP request, this function is only here to satisfy the
// `pack.Server` interface.
func (s *Server) GRPCGatewayClients() []string {
	return []string{}
}

// GRPCGatewayServer indicates the `grpc-gateway` server that the gRPC service is connected to.
func (s *Server) GRPCGatewayServer() string {
	return s.config.GRPCGatewayServer
}

// PreStartCallback is the callback function to trigger right before the server starts running.
func (s *Server) PreStartCallback() func() error {
	return s.preStartCallback
}

// Serve accepts incoming connections on the listener lis, creating a new ServerTransport and
// service goroutine for each. The service goroutines read gRPC requests and then call the
// registered handlers to reply to them. Serve returns when lis.Accept fails with fatal errors.
// lis will be closed when this method returns. Serve will return a non-nil error unless Stop or
// GracefulStop is called.
func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", s.Addr())
	if err != nil {
		return err
	}
	s.listener = lis

	return s.server.Serve(s.listener)
}

// Type indicates if the server is gRPC server.
func (s *Server) Type() string {
	return "gRPC"
}

// TracerProviderShutdownHandler is a function that shuts down the tracer's exporter/provider before
// the gRPC server is gracefully shut down.
func (s *Server) TracerProviderShutdownHandler() error {
	return s.config.TracerProviderShutdownHandler()
}

func defaultServerConfig(c *ServerConfig) {
	if c.KeepAlive.EnforcementPolicy.MinTime == 0 {
		c.KeepAlive.EnforcementPolicy.MinTime = 5 * time.Second
	}

	if c.KeepAlive.ServerParameters.MaxConnectionAge == 0 {
		c.KeepAlive.ServerParameters.MaxConnectionAge = 30 * time.Second
	}

	if c.KeepAlive.ServerParameters.MaxConnectionIdle == 0 {
		c.KeepAlive.ServerParameters.MaxConnectionIdle = 15 * time.Second
	}

	if c.KeepAlive.ServerParameters.Time == 0 {
		c.KeepAlive.ServerParameters.Time = 5 * time.Second
	}

	if c.KeepAlive.ServerParameters.Timeout == 0 {
		c.KeepAlive.ServerParameters.Timeout = 1 * time.Second
	}
}

func loggingInterceptor(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field) {
	if ce := ctxzap.Extract(ctx).Check(level, msg); ce != nil {
		ce.Write(
			zap.Error(err),
			zap.String("grpc.code", code.String()),
			zap.String("trace_id", tracer.GetTraceIDFromContext(ctx)),
			duration,
		)
	}
}
