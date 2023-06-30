# Go-Sdk

The SDK contains commonly used golang packages like HTTP middlewares, GRPC libraries, Excel Sheet Reader/Writer, QR code generator, SQL database building block, Logger & Tracer by different microservices.

Here are some example use cases:

## REST Api server Example

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sdkhttp "github.com/Raj63/go-sdk/http"
	sdkgin "github.com/Raj63/go-sdk/http/gin"
	"github.com/Raj63/go-sdk/logger"
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

```

## GRPC server Example

```go
    package main

import (
	"log"

	sdkgrpc "github.com/Raj63/go-sdk/grpc"
	"github.com/Raj63/go-sdk/logger"
	"google.golang.org/grpc"
)

func main() {
	_logger := logger.NewLogger()

	// Setup the gRPC/HTTP server.
	server, err := sdkgrpc.NewServer(
		&sdkgrpc.ServerConfig{},
		_logger,
		func() error {
			return nil
		},
		[]grpc.UnaryServerInterceptor{}...,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Register the gRPC server implementation.
	// api.RegisterAuthServiceServer(
	// 	server.GRPCServer(),
	// 	&pkggrpc.Server{},
	// )

	// finally serve the server
	if err := server.Serve(); err != nil {
		_logger.Fatal(err)
	}

}

```

## Excel Read/Write Example

```go
package main

import (
	"fmt"

	"github.com/Raj63/go-sdk/excel"
	"github.com/Raj63/go-sdk/excel/excelize"
	"github.com/Raj63/go-sdk/logger"
)

func main() {
	ReadSample()
	ReadSheet()
	ReadRow()
}

// read samples of each sheets from a excel 
func ReadSample() {
	logger := logger.NewLogger()
	excel := excelize.NewExcel(excel.FileInfo{FilePath: "some-excel-file.xlsx"}, logger)

	data, err := excel.ReadSample()
	if err != nil {
		fmt.Printf("Error reading: %v", err)
		return
	}

	for _, sheet := range data.Sheets {
		fmt.Printf("Sheet: %s\n", sheet.Name)
		for _, row := range sheet.DataRow {
			for _, data := range row.Data {
				fmt.Printf("%s: %s\n", data.Header, data.Value)
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func ReadSheet() {
	logger := logger.NewLogger()
	excel := excelize.NewExcel(excel.FileInfo{FilePath: "some-excel-file.xlsx"}, logger)

	sheet, err := excel.ReadSheet("Sheet 1")
	if err != nil {
		fmt.Printf("Error reading: %v", err)
		return
	}

	fmt.Printf("Sheet: %s\n", sheet.Name)
	for _, row := range sheet.DataRow {
		for _, data := range row.Data {
			fmt.Printf("%s: %s\n", data.Header, data.Value)
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func ReadRow() {
	logger := logger.NewLogger()
	exl := excelize.NewExcel(excel.FileInfo{FilePath: "some-excel-file.xlsx"}, logger)

	executor := func(row excel.DataRow) error {
		for _, data := range row.Data {
			fmt.Printf("%s: %s\n", data.Header, data.Value)
		}
		fmt.Println("")
		return nil
	}
	err := exl.ReadRow("Sheet 1", executor)
	if err != nil {
		fmt.Printf("Error reading: %v", err)
		return
	}
}

```

## QR Code Generation Example

```go
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	myqr "github.com/Raj63/go-sdk/qrcode"
	"github.com/Raj63/go-sdk/qrcode/yeqown"

	"github.com/pkg/errors"
)

func main() {
	TestQrCodeWithWriter()
	TestQrCode()
	TestQrCodeN(10)
}

func TestQrCodeWithWriter() {
	//setup qr config
	qrConfig := myqr.DefaultQrConfig()
	qrConfig.LogoPath = "../../assets/images/logo-100.jpeg"
	qrConfig.Shape = myqr.CircleShape

	qr, err := yeqown.NewQrcode(qrConfig)
	if err != nil {
		log.Fatal(err)
	}

	file, err := NewFile("./generated/qrcode-with-url-logo.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	err = qr.CreateWithWriter(myqr.ProductInfo{Host: "https://amritmahotsav.nic.in/writereaddata/Portal/Images/89/290537667_353719030246146_3235967459984024871_n.jpg"}, file)
	if err != nil {
		log.Fatal(err)
	}
}

func TestQrCode() {
	//setup qr config
	qrConfig := myqr.DefaultQrConfig()
	qrConfig.LogoPath = "../../assets/images/logo-100.jpeg"
	qrConfig.Shape = myqr.RectangleShape

	qr, err := yeqown.NewQrcode(qrConfig)
	if err != nil {
		log.Fatal(err)
	}

	file, err := NewFile("./generated/qrcode-with-url-logo.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	result, err := qr.Create(myqr.ProductInfo{Text: `
Name:	Medha Biswas
Guardian Name:	Rajesh Biswas
Emergency Contact:	9902260563
Blood Group:	A+
Address:	1-108/1, Hyderabad
School:	Mayalaxmi International School
School Contact:	8702711233
Class:	LKG
Roll Number:	M23CLR3`})
	if err != nil {
		log.Fatal(err)
	}

	file.Write(result.QRCode)
}

func TestQrCodeN(len int) {
	//setup qr config
	qrConfig := myqr.DefaultQrConfig()
	qrConfig.LogoPath = "../../assets/images/logo-100.jpeg"
	qrConfig.Shape = myqr.CircleShape

	qr, err := yeqown.NewQrcode(qrConfig)
	if err != nil {
		log.Fatal(err)
	}

	products := make([]myqr.ProductInfo, 0)
	for i := 0; i < len; i++ {
		products = append(products, myqr.ProductInfo{
			Text: fmt.Sprintf("text-%d", i),
		})
	}

	result, err := qr.CreateN(products)
	if err != nil {
		log.Fatal(err)
	}

	for _, eachQR := range result {
		{
			file, err := NewFile(fmt.Sprintf("./generated/qrcode-range-%s.jpeg", eachQR.ProductInfo.Text))
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			file.Write(eachQR.QRCode)
		}
	}
}

func NewFile(filename string) (io.WriteCloser, error) {
	if _, err := os.Stat(filename); err != nil && os.IsExist(err) {
		// custom path got: "file exists"
		return nil, fmt.Errorf("could not find path: %s", filename)
	}

	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "create file failed")
	}

	return fd, nil
}

```