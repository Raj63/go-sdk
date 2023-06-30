package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Raj63/go-sdk/errors"
)

// Download is a helper function to download resources for the specified path from http or local
func Download(path string) (io.Reader, error) {
	if path == "" {
		return nil, errors.NewAppErrorWithType(errors.InputEmpty)
	}

	var reader io.Reader
	var err error
	if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "https") {
		// TODO: download the file and return it in the image format
	} else {
		reader, err = os.Open(filepath.Clean(path))
		if err != nil {
			return nil, fmt.Errorf("could not open local file(%s), error=%v", path, err)
		}
	}
	return reader, nil
}
