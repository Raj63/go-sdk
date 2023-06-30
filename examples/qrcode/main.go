package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	myqr "github.com/Raj63/go-sdk/qrcode"
	"github.com/Raj63/go-sdk/qrcode/yeqown"

	"github.com/pkg/errors"
)

func main() {
	testQrCodeWithWriter()
	testQrCode()
	testQrCodeN(10)
}

func testQrCodeWithWriter() {
	// setup qr config
	qrConfig := myqr.DefaultQRConfig()
	qrConfig.LogoPath = "../../assets/images/logo-100.jpeg"
	qrConfig.Shape = myqr.CircleShape

	qr, err := yeqown.NewQrcode(qrConfig)
	if err != nil {
		log.Fatal(err)
	}

	file, err := newFile("./generated/qrcode-with-url-logo.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	err = qr.CreateWithWriter(myqr.ProductInfo{Host: "https://amritmahotsav.nic.in/writereaddata/Portal/Images/89/290537667_353719030246146_3235967459984024871_n.jpg"}, file)
	if err != nil {
		log.Fatal(err)
	}
}

func testQrCode() {
	// setup qr config
	qrConfig := myqr.DefaultQRConfig()
	qrConfig.LogoPath = "../../assets/images/logo-100.jpeg"
	qrConfig.Shape = myqr.RectangleShape

	qr, err := yeqown.NewQrcode(qrConfig)
	if err != nil {
		log.Fatal(err)
	}

	file, err := newFile("./generated/qrcode-with-url-logo.jpeg")
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
		fmt.Println(err)
	}

	_, err = file.Write(result.QRCode)
	if err != nil {
		fmt.Println(err)
	}
}

func testQrCodeN(len int) {
	// setup qr config
	qrConfig := myqr.DefaultQRConfig()
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
			file, err := newFile(fmt.Sprintf("./generated/qrcode-range-%s.jpeg", eachQR.ProductInfo.Text))
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err = file.Close(); err != nil {
					fmt.Println(err)
				}
			}()
			_, err = file.Write(eachQR.QRCode)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func newFile(filename string) (io.WriteCloser, error) {
	if _, err := os.Stat(filename); err != nil && os.IsExist(err) {
		// custom path got: "file exists"
		return nil, fmt.Errorf("could not find path: %s", filename)
	}

	fd, err := os.OpenFile(filepath.Clean(filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "create file failed")
	}

	return fd, nil
}
