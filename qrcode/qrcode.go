package qrcode

import (
	"fmt"
	"io"
	"net/url"
)

// QRCode specifies the contraints related to qrcode generation
type QRCode interface {
	Create(product ProductInfo) (QRResult, error)
	CreateN(product []ProductInfo) ([]QRResult, error)
	CreateWithWriter(product ProductInfo, writer io.WriteCloser) error
}

// ProductInfo represents a product
type ProductInfo struct {
	Host string
	Text string
}

// QRResult represent the result with qrcode
type QRResult struct {
	ProductInfo ProductInfo
	QRCode      []byte
}

// Content returns the value from ProductInfo for which qrcode will be generated
func (p ProductInfo) Content() (string, error) {
	text := p.Text
	if p.Host != "" {
		resourceURL, err := url.Parse(p.Host)
		if err != nil {
			return "", fmt.Errorf("invalid Host url %q", p.Host)
		}
		text = resourceURL.String()
	}
	return text, nil
}
