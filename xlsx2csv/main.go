package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

var xlsxPath = flag.String("f", "", "Path to an XLSX file")
var sheetIndex = flag.Int("i", 0, "Index of sheet to convert, zero based")
var sheetName = flag.String("s", "", "Name of sheet to convert")
var delimiter = flag.String("d", ";", "Delimiter to use between fields")
var startrow = flag.Int("r", 1 , "Start xlsx read from line nr, default 1")

type outputer func(s string)

func generateCSVFromXLSXFile(excelFileName string, sheetIndex int, sheetName string, outputf outputer) error {
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("This XLSX file contains no sheets.")
	case sheetIndex >= sheetLen:
		return fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
	}
	sheet := xlFile.Sheets[sheetIndex]
	// or like to use sheet name ?
	if sheetName != ""  {
		sheet2, ok := xlFile.Sheet[sheetName]
		if ok != true {
			return errors.New("This XLSX file contains not named sheet.")
		}
		sheet=sheet2
	}
	var rownr int = 0
	for _, row := range sheet.Rows {
		rownr++
		if rownr+1 < *startrow { continue }
		var vals []string
		if row != nil {
			for _, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil {
					vals = append(vals, err.Error())
				}
				vals = append(vals, fmt.Sprintf("%q", str))
			}
			outputf(strings.Join(vals, *delimiter) + "\n")
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if len(os.Args) < 3 {
		flag.PrintDefaults()
		return
	}
	flag.Parse()
	printer := func(s string) { fmt.Printf("%s", s) }
	if err := generateCSVFromXLSXFile(*xlsxPath, *sheetIndex, *sheetName, printer); err != nil {
		fmt.Println(err)
	}
}
