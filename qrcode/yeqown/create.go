package yeqown

import (
	"fmt"

	"github.com/Raj63/go-sdk/qrcode"
)

// Create implements qrcode.QRCode, it create and returns the QRResult
func (y *yeqownQrcode) Create(product qrcode.ProductInfo) (qrcode.QRResult, error) {
	byteArrayWriteCloser := qrcode.NewByteArrayWriteCloser()
	err := y.CreateWithWriter(product, byteArrayWriteCloser)
	if err != nil {
		return qrcode.QRResult{}, fmt.Errorf("error creating Qrcode with ByteArrayWriterCloser: %w", err)
	}
	return qrcode.QRResult{ProductInfo: product, QRCode: byteArrayWriteCloser.Bytes()}, nil
}
