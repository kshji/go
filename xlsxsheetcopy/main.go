// go run main.go -f my.xlsx -i 0 -n Copied,Copied2,101,102 --debug 1 

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/tealeg/xlsx"
)

var xlsxPath = flag.String("f", "", "XLSX input file")
var xlsxOut = flag.String("o", "", "XLSX output file")
var sheetIndex = flag.Int("i", 0, "Index of sheet to copy, zero based")
var sheetName = flag.String("s", "", "Name of sheet to copy")
var sheetNew = flag.String("n", "Sheet New", "List of New Sheet Names comma separated")
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

	// Sheet name to duplicate
	sheet1Name := sheet.Name
	if *debug>0 { fmt.Println("Duplicate SheetName:",sheet1Name ) }
	// duplicate org sheet

	// name of new sheets is comma separeted list sheet names
	for  _, sheetname := range strings.Split(sheetNameNew,",") {
		if *debug>0 { fmt.Println("New Sheet:",sheetname ) }
		xlFile.AppendSheet(*sheet,sheetname)
	}
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
