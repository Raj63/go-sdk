// Package yeqown is a wrapper around library Ref: https://github.com/yeqown/go-qrcode
package yeqown

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path"

	"github.com/Raj63/go-sdk/errors"

	"github.com/Raj63/go-sdk/file"
	"github.com/Raj63/go-sdk/qrcode"

	"github.com/yeqown/go-qrcode/writer/standard"
)

type yeqownQrcode struct {
	options []standard.ImageOption
}

// NewQrcode creates aninstance of the yeqownQrcode with the given QrConfig
func NewQrcode(qrConfig qrcode.QRConfig) (qrcode.QRCode, error) {
	options := make([]standard.ImageOption, 0)
	options = append(options, standard.WithBgColorRGBHex(qrConfig.BgHexColor))
	options = append(options, standard.WithFgColorRGBHex(qrConfig.FgHexColor))

	logo, err := downloadLogoImage(qrConfig.LogoPath)
	if err != nil {
		return nil, fmt.Errorf("error downloading logo Img: %w", err)
	}
	options = append(options, standard.WithLogoImage(logo))
	options = append(options, standard.WithQRWidth(uint8(qrConfig.Width)))

	// if qrConfig.Shape == qrcode.CircleShape {
	// 	options = append(options, standard.WithCircleShape())
	// }

	return &yeqownQrcode{
		options: options,
	}, nil
}

func downloadLogoImage(path string) (image.Image, error) {
	reader, err := file.Download(path)
	if errors.ErrorType(err) == errors.InputEmpty {
		return nil, nil
	}
	return decodeImage(path, reader)
}

func decodeImage(filePath string, reader io.Reader) (image.Image, error) {
	switch path.Ext(filePath) {
	case ".jpeg", ".jpg":
		img, err := jpeg.Decode(reader)
		if err != nil {
			return nil, fmt.Errorf("could not open file(%s), error=%v", filePath, err)
		}
		return img, nil

	case ".png":
		img, err := png.Decode(reader)
		if err != nil {
			return nil, fmt.Errorf("could not open file(%s), error=%v", filePath, err)
		}
		return img, nil

	}

	return nil, fmt.Errorf("file not supported, use Jpeg or PNG format")
}
