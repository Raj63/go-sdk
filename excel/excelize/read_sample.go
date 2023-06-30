package excelize

import (
	"fmt"

	"github.com/Raj63/go-sdk/excel"
)

const (
	// SampleSize indicates the size of the sample rows
	SampleSize = 3
)

// ReadSample implements excel.Excel
func (e *excelize) ReadSample() (*excel.Data, error) {
	f, closer, err := e.GetFile()
	if err != nil {
		return nil, err
	}
	defer closer()

	sheets := make([]excel.Sheet, 0)
	for _, sheetName := range f.GetSheetList() {
		// Get all the rows in the sheetName.
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("read rows: %w", err)
		}
		sheet := excel.Sheet{
			Name:    sheetName,
			DataRow: make([]excel.DataRow, 0),
		}
		headerMaps := make(map[int]string, 0)
		for rowIndex, row := range rows {
			// sampling rows with SampleSize
			if rowIndex == SampleSize {
				break
			}
			dataRow := excel.DataRow{
				Data: make([]excel.KeyValue, 0),
			}
			for colIndex, colCell := range row {
				if rowIndex == 0 {
					headerMaps[colIndex] = colCell
				} else {
					dataRow.Data = append(dataRow.Data, excel.KeyValue{
						Header: headerMaps[colIndex],
						Value:  colCell,
					})
				}
			}
			sheet.DataRow = append(sheet.DataRow, dataRow)
		}
		sheets = append(sheets, sheet)
	}
	return &excel.Data{
		Name:   e.name,
		Sheets: sheets,
	}, nil
}
