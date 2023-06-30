package excelize

import (
	"fmt"

	"github.com/Raj63/go-sdk/errors"
	"github.com/Raj63/go-sdk/excel"
)

// ReadRow implements excel.Excel
func (e *excelize) ReadRow(inputSheetName string, exec func(row excel.DataRow) error) error {
	f, closer, err := e.GetFile()
	if err != nil {
		return err
	}
	defer closer()

	for _, sheetName := range f.GetSheetList() {
		if sheetName == inputSheetName {
			// Get all the rows in the sheetName.
			rows, err := f.GetRows(sheetName)
			if err != nil {
				return fmt.Errorf("read rows: %w", err)
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
				// invoke the execute function
				if err := exec(dataRow); err != nil {
					return fmt.Errorf("failed to execute func: %w", err)
				}
			}
			return nil
		}
	}
	return errors.NewAppErrorWithType(errors.NotFound)
}
