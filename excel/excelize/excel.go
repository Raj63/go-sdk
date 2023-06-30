package excelize

import (
	"fmt"
	"io"

	"github.com/Raj63/go-sdk/errors"
	"github.com/Raj63/go-sdk/excel"
	"github.com/Raj63/go-sdk/file"
	"github.com/Raj63/go-sdk/logger"
	exc "github.com/xuri/excelize/v2"
)

type excelize struct {
	name       string
	filePath   string
	logger     *logger.Logger
	fileReader io.Reader
}

// NewExcel returns the excelize object which implements the excel
func NewExcel(input excel.FileInfo, logger *logger.Logger) excel.Excel {
	return &excelize{
		logger:     logger,
		filePath:   input.FilePath,
		fileReader: input.FileReader,
	}
}

func (e *excelize) GetFile() (*exc.File, func(), error) {
	if e.fileReader == nil && e.filePath != "" {
		// download the file
		reader, err := file.Download(e.filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("download file: %w", err)
		}
		e.fileReader = reader
	}

	if e.fileReader == nil {
		return nil, nil, errors.NewAppErrorWithType(errors.InputEmpty)
	}

	f, err := exc.OpenReader(e.fileReader)
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}

	return f, func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			e.logger.Error("close spreadsheet file: %w", err)
		}
	}, nil
}
