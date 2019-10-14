package main

import (
	"fmt"
	"os"

	"encoding/csv"
	"io"
	"strconv"
	//"strings"
	//"math"

	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
	"unicode/utf8"
)

// SheetNamesTemplate define name's for new created sheets.
var SheetNamesTemplate = "Sheet %d"

func main() {
	initCommandLine(os.Args)
}

func writeAllSheets(xlFile *xlsx.File, dataFiles []string, sheetNames []string, exampleRowNumber int) (err error) {

	for i, dataFileName := range dataFiles {

		sheet, err := getSheet(xlFile, sheetNames, i)
		if err != nil {
			return err
		}

		var exampleRow *xlsx.Row
		if exampleRowNumber != 0 && exampleRowNumber <= len(sheet.Rows) {
			// example row counting from 1
			exampleRow = sheet.Rows[exampleRowNumber-1]

			// remove example row
			sheet.Rows = append(sheet.Rows[:exampleRowNumber-1], sheet.Rows[exampleRowNumber:]...)
		}

		err = writeSheet(dataFileName, sheet, exampleRow)

		if err != nil {
			return err
		}
	}

	return nil
}

func getSheet(xlFile *xlsx.File, sheetNames []string, i int) (sheet *xlsx.Sheet, err error) {

	var sheetName string
	if len(sheetNames) > i {
		sheetName = sheetNames[i]
	} else {
		sheetName = fmt.Sprintf(SheetNamesTemplate, i+1)
	}

	sheet, ok := xlFile.Sheet[sheetName]
	if ok != true {
		sheet, err = xlFile.AddSheet(sheetName)

		if err != nil {
			return nil, err
		}
	}
	return sheet, nil
}

func writeSheet(dataFileName string, sheet *xlsx.Sheet, exampleRow *xlsx.Row) error {

	data, err := getCsvData(dataFileName)
	//data.Comma=';'
	//data.Comma=colSep
	//fmt.Printf("%c", []rune(colSep)[1])
	//data.Comma=colSep[:1]
	//r, _ := utf8.DecodeRuneInString(colSep)
	r, _ := utf8.DecodeRuneInString(myParam.colsep)
	data.Comma=r

	if myParam.writeheaderlines == 0 && myParam.startrow < (myParam.headerlines+1) {
		// startline have to be headerlines + 1
		myParam.startrow=myParam.headerlines+1

		if (myParam.debug>0) { fmt.Println("Startrow have to be >= headerlines+1 = ",myParam.startrow) }
		}
	if (myParam.debug>0) {
		//fmt.Printf("Comma:%#v",data.Comma)
		//fmt.Println()
		fmt.Printf("Comma:%q", r)
		fmt.Println()
		fmt.Printf("Comma:%s", myParam.colsep)
		fmt.Println()
		fmt.Println("Comma:", myParam.colsep)
		//fmt.Println('Comma:', myParam.output)
		fmt.Println( myParam.output)
		//fmt.Printf('Comma:%s', param.output)
		fmt.Println()
		fmt.Println("startrow",myParam.startrow)
		fmt.Println("headerlines",myParam.headerlines)
		fmt.Println("writeheaderlines",myParam.writeheaderlines)
		}
	//data.Comma=rune(colSep[0:1])
	//r := csv.NewReader(strings.NewReader(in))
	//r.Comma = ';'
	//r.Comment = '#'

	

	if err != nil {
		return err
	}

	var i int
	var writeline int
	writeline=0

	for {
		record, err := data.Read()
		i++
		writeline=0
		if myParam.writeheaderlines == 1 { writeline=1 }
		if  myParam.debug>0  && record != nil {
			fmt.Println("rivi:",i,record,writeline,myParam.startrow,myParam.headerlines);
			}

		if err == io.EOF || record == nil {
			break
		} else if err != nil {
			return err
		}

		//if i > 5000 {
		//	break
		//}

		//if i%500 == 0 {
		//	fmt.Println(i)
		//}



		if myParam.verbose>0 { fmt.Printf("\r%08d",i) }

		if writeline>0 || i>= myParam.startrow || i> myParam.headerlines { writeRowToXls(sheet, record, exampleRow,i)  }
	}

	if myParam.verbose>0 {
		fmt.Println()
		fmt.Println()
		}

	return nil
}

func buildXls(c *cli.Context, p *params) (err error) {

	var xlFile *xlsx.File
	if p.xlsxTemplate == "" {
		xlFile = xlsx.NewFile()
	} else {
		xlFile, err = xlsx.OpenFile(p.xlsxTemplate)
		if err != nil {
			return err
		}
	}

	writeAllSheets(xlFile, p.input, p.sheets, p.row)

	return xlFile.Save(p.output)
}


func writeRowToXls(sheet *xlsx.Sheet, record []string, exampleRow *xlsx.Row, rownr int) {

	var row *xlsx.Row
	var cell *xlsx.Cell
	//var celltype *xlsx.CellType

	row = sheet.AddRow()
	//row.WriteSlice( &record , -1)

	var cellsLen int
	if exampleRow != nil {
		cellsLen = len(exampleRow.Cells)
	}

	// write cells
	for k, v := range record {
		cell = row.AddCell()
		if myParam.debug>0 && rownr==1 { fmt.Println("Col:",k,v) }
		if exampleRow != nil && cellsLen > k {  // example row, use it
			writeCell(cell, exampleRow, k, v, rownr)
		} else  {  // no example row, so try to quess coltype, setCellValue do it
				if []rune(v)[0] == '='  {   // input data include =formula syntax
					cell.SetFormula(v[1:])
				} else {
					setCellValue(cell, v)
				}
			}
	} //for
}

// Write Cell using ExampleRow
func writeCell(cell *xlsx.Cell, exampleRow *xlsx.Row, colNr int, colString string, rownr int) {

	var formula string

	formula = ""
	colStr := colString
	// get formula , it's empty string if no formula
        formula = exampleRow.Cells[colNr].Formula()
	celltype := exampleRow.Cells[colNr].Type()
	cStyle := exampleRow.Cells[colNr].GetStyle()

	if []rune(colStr)[0] == '='  {   // input data include =formula syntax, use it, not examplerow
		formula = colStr[1:]  // remove 1st =
		colStr=""  // have to be empty, then formula has to be set
		}
	if rownr <= myParam.headerlines { // headerlines, allways default
		formula=""
		celltype=0
		if myParam.debug>0  { fmt.Println("Headerline:",colNr,colString) }
		}
	if formula != ""  && (colStr=="" || colStr=="-" ) { // set Formula, not value
		cell.SetFormula(formula)
		cFormat := exampleRow.Cells[colNr].GetNumberFormat()
		cell.SetFormat(cFormat)
		cell.SetStyle(cStyle)
		if (myParam.debug>0  ) {
			fmt.Printf(" - Formula:%s",formula)
			fmt.Printf(" Format %s",cFormat)
			fmt.Println()
			}
		return
		}
	if myParam.debug>0 && rownr==1 { fmt.Println("Type:",celltype) }
	// - allways 2 = xlsx.CellTypeNumeric ????
	switch celltype {
		case xlsx.CellTypeNumeric:
			if myParam.debug>0 && rownr==1  {  fmt.Println("CellTypeNumeric")   }
			floatVal, err := strconv.ParseFloat(colStr, 64)
			if err != nil {
				setCellValue(cell, colStr)
				//fmt.Printf("(%d,%d) is not a valid number, value: %s", rownr,k, v)
				//fmt.Println()
			} else {
				cFormat := exampleRow.Cells[colNr].GetNumberFormat()
				//cell.SetFloatWithFormat(floatVal, "0.00")
				if myParam.debug>0 && rownr==1 {
					fmt.Printf("Format %s",cFormat)
					fmt.Println()
					}
				//cell.SetFloatWithFormat(floatVal, "0.00")
				cell.SetFloatWithFormat(floatVal, cFormat)
				//setCellValue(cell, colStr)
				//cell.SetFormat(cFormat)
				if myParam.debug>0 && rownr==1 {
					fmt.Printf("Floatf %f",floatVal)
					fmt.Println("Value Set")
					}
				}
		default:
			if myParam.debug>0 && rownr==1  {  fmt.Println("celltype Default")   }
			setCellValue(cell, colStr)
		}

	cell.SetStyle(cStyle)
	celltype =  exampleRow.Cells[colNr].Type()
	if myParam.debug>0 && rownr==1   {
		fmt.Println("After set, Type:",celltype)
		}

}

// setCellValue set data in correct format.
func setCellValue(cell *xlsx.Cell, v string) {

	intVal, err := strconv.Atoi(v)
	if err == nil {
		if intVal < 100000000000 { // Long numbers are displayed incorrectly in Excel
			cell.SetInt(intVal)
			return
		}
		cell.Value = v
		return
	}

	floatVal, err := strconv.ParseFloat(v, 64)
	if err == nil {
		cell.SetFloat(floatVal)
		return
	}
	cell.Value = v
}

// getCsvData read's data from CSV file.
func getCsvData(dataFileName string) (*csv.Reader, error) {

	dataFile, err := os.Open(dataFileName)
	if err != nil {
		return nil, cli.NewExitError("Problem with reading data from "+dataFileName, 11)
	}

	return csv.NewReader(dataFile), nil
}
