package qrcode

// QRConfig represnts the configuration to generate a QR Code
type QRConfig struct {
	// Validates logo size as follows Width >= 2*logoWidth && qrHeight >= 2*logoHeight
	// Instead of default expression Width >= 5*logoWidth && qrHeight >= 5*logoHeight
	LogoPath   string
	BgHexColor string
	FgHexColor string
	Width      int
	Shape      ShapeType
}

// ShapeType specifies the type of shape to generate a qr code
type ShapeType string

const (
	// RectangleShape will generate the QR Code with rectangular shape
	RectangleShape ShapeType = "rectangle"
	// CircleShape will generate the QR Code with circular shape
	CircleShape ShapeType = "circle"
)

// DefaultQRConfig returns a default QRConfig
func DefaultQRConfig() QRConfig {
	return QRConfig{
		BgHexColor: "#FFFFFF",
		FgHexColor: "#000000",
	}
}
