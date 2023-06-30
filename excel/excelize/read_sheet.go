package excelize

import (
	"fmt"

	"github.com/Raj63/go-sdk/errors"
	"github.com/Raj63/go-sdk/excel"
)

// ReadSheet implements excel.Excel
func (e *excelize) ReadSheet(inputSheetName string) (*excel.Sheet, error) {
	f, closer, err := e.GetFile()
	if err != nil {
		return nil, err
	}
	defer closer()

	for _, sheetName := range f.GetSheetList() {
		if sheetName == inputSheetName {
			// Get all the rows in the sheetName.
			rows, err := f.GetRows(sheetName)
			if err != nil {
				return nil, fmt.Errorf("read rows: %w", err)
			}
			sheet := &excel.Sheet{
				Name:    sheetName,
				DataRow: make([]excel.DataRow, 0),
			}
			headerMaps := make(map[int]string, 0)
			for rowIndex, row := range rows {
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
			return sheet, nil
		}
	}
	return nil, errors.NewAppErrorWithType(errors.NotFound)
}
