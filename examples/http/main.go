package main

import (
	"log"
	"net/http"
	"time"

	sdkhttp "github.com/Raj63/go-sdk/http"
	sdkgin "github.com/Raj63/go-sdk/http/gin"
	"github.com/Raj63/go-sdk/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize the router
	router := gin.Default()

	_logger := logger.NewLogger()

	// make sure the middleware configs are passed for HTTP server
	err := sdkgin.AddBasicHandlers(router, &sdkgin.MiddlewaresConfig{
		DebugEnabled: true,
		RateLimiterConfig: struct {
			Enabled    bool
			Interval   time.Duration
			BucketSize int
		}{
			Enabled:    true,
			Interval:   time.Second * 1,
			BucketSize: 3,
		},
		CorsOptions: struct {
			Enabled         bool
			AllowOrigins    []string
			AllowMethods    []string
			AllowHeaders    []string
			ExposeHeader    []string
			AllowOriginFunc func(origin string) bool
			MaxAge          time.Duration
		}{
			Enabled:      true,
			AllowOrigins: []string{"http://localhost:8080"},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
			AllowHeaders: []string{"Origin"},
			ExposeHeader: []string{"Content-Length"},
		},
		StaticFilesOptions: struct {
			Enabled    bool
			ServeFiles []struct {
				Prefix                 string
				FilePath               string
				AllowDirectoryIndexing bool
			}
		}{
			Enabled: true,
			ServeFiles: []struct {
				Prefix                 string
				FilePath               string
				AllowDirectoryIndexing bool
			}{
				{
					Prefix:                 "/",
					FilePath:               "assets/static",
					AllowDirectoryIndexing: true,
				},
			},
		},
		PrometheusEnabled: true,
		NewRelicOptions: struct {
			ServiceName string
			LicenseKey  string
		}{
			ServiceName: "ServiceName",
			LicenseKey:  "NewRelicLicenseKey",
		},
	}, _logger)
	if err != nil {
		_logger.Errorf("error setting up HTTP basic middlewares: %v", err)
		log.Fatalln(err)
	}

	// register the application routes here
	// routes.ApplicationV1Router(router)

	httpServer, err := sdkhttp.NewServer(
		&sdkhttp.ServerConfig{
			Name:    "service-name",
			Address: "0.0.0.0:8080",
			GracefulShutdownHandler: func() error {
				// gracefully shutdown database stuff
				return nil
			},
		},
		_logger,
		func() error { return nil },
		router)
	if err != nil {
		_logger.Errorf("error creating server: %v", err)
		log.Fatalln(err)
	}

	// finally serve the server
	if err := httpServer.Serve(); err != nil {
		_logger.Fatal(err)
	}
}
