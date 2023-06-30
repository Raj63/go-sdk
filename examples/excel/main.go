package main

import (
	"fmt"

	"github.com/Raj63/go-sdk/excel"
	"github.com/Raj63/go-sdk/excel/excelize"
	"github.com/Raj63/go-sdk/logger"
)

func main() {
	readSample()
	readSheet()
	readRow()
}

func readSample() {
	logger := logger.NewLogger()
	excel := excelize.NewExcel(excel.FileInfo{FilePath: "../../assets/excel/Mayalaxmi_International_School.xlsx"}, logger)

	data, err := excel.ReadSample()
	if err != nil {
		fmt.Printf("Error reading: %v", err)
		return
	}

	for _, sheet := range data.Sheets {
		fmt.Printf("Sheet: %s\n", sheet.Name)
		for _, row := range sheet.DataRow {
			for _, data := range row.Data {
				fmt.Printf("%s: %s\n", data.Header, data.Value)
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func readSheet() {
	logger := logger.NewLogger()
	excel := excelize.NewExcel(excel.FileInfo{FilePath: "../../assets/excel/Mayalaxmi_International_School.xlsx"}, logger)

	sheet, err := excel.ReadSheet("Sheet 1")
	if err != nil {
		fmt.Printf("Error reading: %v", err)
		return
	}

	fmt.Printf("Sheet: %s\n", sheet.Name)
	for _, row := range sheet.DataRow {
		for _, data := range row.Data {
			fmt.Printf("%s: %s\n", data.Header, data.Value)
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func readRow() {
	logger := logger.NewLogger()
	exl := excelize.NewExcel(excel.FileInfo{FilePath: "../../assets/excel/Mayalaxmi_International_School.xlsx"}, logger)

	executor := func(row excel.DataRow) error {
		for _, data := range row.Data {
			fmt.Printf("%s: %s\n", data.Header, data.Value)
		}
		fmt.Println("")
		return nil
	}
	err := exl.ReadRow("Sheet 1", executor)
	if err != nil {
		fmt.Printf("Error reading: %v", err)
		return
	}
}
