package yeqown

import (
	"fmt"
	"io"

	"github.com/Raj63/go-sdk/qrcode"

	goqrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// CreateWithWriter implements qrcode.QRCode, which creates QR code for given ProductInfo
func (y *yeqownQrcode) CreateWithWriter(product qrcode.ProductInfo, writer io.WriteCloser) error {
	content, err := product.Content()
	if err != nil {
		return fmt.Errorf("invalid input product info: %w", err)
	}
	qrc, err := goqrcode.NewWith(content,
		goqrcode.WithEncodingMode(goqrcode.EncModeByte),
		goqrcode.WithErrorCorrectionLevel(goqrcode.ErrorCorrectionQuart),
	)
	if err != nil {
		return fmt.Errorf("error initializing qrcode: %w", err)
	}

	w := standard.NewWithWriter(
		writer,
		y.options...,
	)

	if err != nil {
		return fmt.Errorf("error creating qrcode: %w", err)
	}

	if err = qrc.Save(w); err != nil {
		return fmt.Errorf("error writing qrcode: %w", err)
	}
	return nil
}
