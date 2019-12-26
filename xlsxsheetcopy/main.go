// go run main.go -f my.xlsx -i 0 -n Copied --debug 1 

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"github.com/tealeg/xlsx"
)

var xlsxPath = flag.String("f", "", "XLSX input file")
var xlsxOut = flag.String("o", "", "XLSX output file")
var sheetIndex = flag.Int("i", 0, "Index of sheet to copy, zero based")
var sheetName = flag.String("s", "", "Name of sheet to copy")
var sheetNew = flag.String("n", "Sheet New", "New Sheet Name")
var debug = flag.Int("debug", 0, "debug 0|1 ")

func copySheet(excelFileName string, outFile string, sheetIndex int, sheetName string, sheetNameNew string) error {

	//outFile := flag.Args()[0]
	xlFile, error := xlsx.OpenFile(excelFileName)
	//newfile := xlsx.NewFile()
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	if *debug>0 { fmt.Println("SheetLen:",sheetLen) }

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

	// Sheet name to copy
	sheet1Name := sheet.Name
	if *debug>0 { fmt.Println("SheetName:",sheet1Name ) }
	// add org sheet
	//newfile.AppendSheet(*sheet,sheet1Name)
	// add copy of org
	//newfile.AppendSheet(*sheet,sheetNameNew)
	xlFile.AppendSheet(*sheet,sheetNameNew)
	//if *debug>0 { fmt.Println("Save:",outFile ) }
	//error = newfile.Save(outFile)
	error = xlFile.Save(excelFileName)
	if error != nil {
                return error
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
	if err := copySheet(*xlsxPath, *xlsxOut, *sheetIndex, *sheetName, *sheetNew); err != nil {
		fmt.Println(err)
	}
}
