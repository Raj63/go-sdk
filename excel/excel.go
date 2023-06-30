package excel

import "io"

// Excel represents the functionalities required to read and write Excel files
type Excel interface {
	ReadSample() (*Data, error)
	ReadSheet(sheetName string) (*Sheet, error)
	ReadRow(sheetName string, exec func(row DataRow) error) error
}

// FileInfo represents the information about a file
type FileInfo struct {
	Name       string
	FilePath   string
	FileReader io.Reader
}

// Data represents the information about excel data
type Data struct {
	Name   string
	Sheets []Sheet
}

// Sheet represents the information about excel sheet
type Sheet struct {
	Name    string
	DataRow []DataRow
}

// DataRow represents the information about excel row
type DataRow struct {
	Data []KeyValue
}

// KeyValue represents the information about each cell in a sheet
type KeyValue struct {
	Header string
	Value  string
}
