/*
https://golang.hotexamples.com/examples/github.com.tealeg.xlsx/Cell/SetStyle/golang-cell-setstyle-method-examples.html
https://golang.hotexamples.com/examples/github.com.tealeg.xlsx/Cell/-/golang-cell-class-examples.html

strings.Replace(col, "_", " ", -1)
*/

package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"unicode/utf8"
	"io"
	"strconv"
	"regexp"
	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
	"github.com/davecgh/go-spew/spew"
	"github.com/kr/pretty"
)

// SheetNamesTemplate define name's for new created sheets.
var SheetNamesTemplate = "Sheet %d"

// ConfigCol struct
type ConfigCol struct {
        colname string
        id int
        fldtype string
        format string
        align string
        font string
        fontsize int
        bold bool
        underline bool
        italic bool
}

func main() {
	initCommandLine(os.Args)
}

func buildXls(c *cli.Context, p *params) (err error) {

	var xlFile *xlsx.File = nil
	var xlFooter *xlsx.File = nil
	if p.xlsxTemplate == "" {
		xlFile = xlsx.NewFile()
	} else {
		xlFile, err = xlsx.OpenFile(p.xlsxTemplate)
		if err != nil {
			return err
		}
	}

	if p.xlsxFooter != "" {
		if p.debug > 0 { fmt.Println("Footer:",p.xlsxFooter) }
		xlFooter, err = xlsx.OpenFile(p.xlsxFooter)
	}
        xlsx.SetDefaultFont(p.fontsize,p.font)

	writeAllSheets(xlFile, p.input, p.sheets, p.row, xlFooter )

	return xlFile.Save(p.output)
}


// 
func writeAllSheets(xlFile *xlsx.File, dataFiles []string, sheetNames []string, exampleRowNumber int, xlFooter *xlsx.File) (err error) {

	if myParam.debug>0 { fmt.Println("Set DefaultFont:",myParam.fontsize,myParam.font) }


	var sheetFooter *xlsx.Sheet = nil
	var xlRow *xlsx.Row
	var exampleRow *xlsx.Row

	for i, dataFileName := range dataFiles {

		sheet, err := getSheet(xlFile, sheetNames, i)
		if err != nil {
			return err
		}
		if xlFooter != nil {
			if myParam.debug>0 { }
			sheetFooter, err = getSheet(xlFooter, sheetNames, i)
			if err != nil {
				return err
			}
                }

		if exampleRowNumber != 0 && exampleRowNumber <= len(sheet.Rows) {
			// example row counting from 1
			exampleRow = sheet.Rows[exampleRowNumber-1]
			if myParam.debug>0 { fmt.Println("Template row dbg:",exampleRow) }
			// kesken
			//parseExampleSheet(sheet, exampleRowNumber)

			sheet.Rows = append(sheet.Rows[:exampleRowNumber-1], sheet.Rows[exampleRowNumber:]...)
		}
		if exampleRowNumber == 0  { // if need some init ...
			}

		_, err = writeSheet(dataFileName, sheet, exampleRow)
		if err != nil {
			fmt.Println("writeSheet, error end:",err)
			return err
		}

		// if Footer defined then append it
		if sheetFooter != nil {
			if  myParam.debug>0 { fmt.Println("Footer sheet:",sheetFooter) }
			for _, xlsxrow  := range sheetFooter.Rows {
				xlRow = sheet.AddRow()
				*xlRow=*xlsxrow
			}
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


// kesken
// loop Example Sheet and expand variables if exists
func parseExampleSheet(sheet *xlsx.Sheet, exampleRow int) {
	for rownr, row := range sheet.Rows {
		if rownr+1 >= exampleRow { continue }
		//var vals []string
		if row != nil {
			for cellnr, cell := range row.Cells {
				value, err := cell.FormattedValue()
				if err != nil {
					//vals = append(vals, err.Error())
				}
				if myParam.debug>0 { fmt.Println("  ex:",rownr,cellnr,value) }
				//vals = append(vals, fmt.Sprintf("%q", str))
			}
			//outputf(strings.Join(vals, *delimiter) + "\n")
		}
	}
}

// some global values from cols
const MaxCols = 10000
var colsWidth [MaxCols]float64
var colsType [MaxCols]xlsx.CellType

// writeSheet
func writeSheet(dataFileName string, sheet *xlsx.Sheet, exampleRow *xlsx.Row) (int, error) {

	if (myParam.debug>0) { fmt.Println("writeSheet BEGIN") }

	data, err := getCsvData(dataFileName)
	if err != nil {
		return 0, err
	}

	r, _ := utf8.DecodeRuneInString(myParam.colsep)
	data.Comma=r  // set for csv reader


	if myParam.writeheaderlines == 0 && myParam.startrow < (myParam.headerlines+1) {
		// startline have to be headerlines + 1
		myParam.startrow=myParam.headerlines+1

		if (myParam.debug>0) { fmt.Println("Startrow have to be >= headerlines+1 = ",myParam.startrow) }
		}
	if (myParam.debug>0) {
		fmt.Printf("Comma:%q", r)
		fmt.Printf("Comma:%s", myParam.colsep)
		fmt.Println("Comma:", myParam.colsep)
		fmt.Println( myParam.output)
		fmt.Println("startrow",myParam.startrow)
		fmt.Println("headerlines",myParam.headerlines)
		fmt.Println("writeheaderlines",myParam.writeheaderlines)
		fmt.Println("config",myParam.config)
		fmt.Println("font",myParam.font)
		fmt.Println("fontsize",myParam.fontsize)
		}

	// set default font - currently xlsx v2 not work ...
	if myParam.debug>0 { fmt.Println("Set DefaultFont:",myParam.fontsize,myParam.font) }
	xlsx.SetDefaultFont(myParam.fontsize,myParam.font)

	// setup default values for Cols, if we are using configfile

	var i int
	var writeline int
	var colnames []string
	var coltypes []string
	var cols []string
	var dataline  int
	var sheetMaxCols = 0 // read from csv
	dataline=(-1)  // number of datalines after headerline
	writeline=0


	// read csv lines 
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
			return 0, err
		}


		if (i == 1 ) { // 1st csv line, 
			sheetMaxCols=len(record) // set number of Cols from 1st line
			cols=record // default use values even it not colnames , some setup need col vector
			if  myParam.debug>0  {
				fmt.Println("1stline:",i,record,sheetMaxCols,myParam.headerlines);
			}


			if  myParam.headerlines > 0 && i ==  1 { // 1st line and colnames
				// save colnames
				colnames=record
				cols=colnames
				dataline=0  // 1st dataline has next
				// check, if colnames include type informatin [n|d|i], add type sto the array coltypes
				setColsType(cols, &coltypes)
				if myParam.debug>0 { fmt.Println("Colnames:",cols ) }
				if myParam.debug>0 { fmt.Println("Coltypes:",coltypes ) }
				}
			if  myParam.headerlines > 0 && i == 1 && myParam.config != "" { // 1st line look config
				if  myParam.debug>0  {
					fmt.Println("1stline, use config:",myParam.config,record);
				}
				// 1st line process, in this point we will check config for sheet
				useConfig(sheet)
			}
		}

		if myParam.verbose>0 { fmt.Printf("\r%08d",i) }

		if writeline>0 || i>= myParam.startrow || i> myParam.headerlines {
			dataline++
			if dataline == 1 {  // 1st dataline, some special ?
				}
			writeRowToXls(sheet, record, exampleRow,i,colnames, coltypes, sheetMaxCols, dataline)
			if dataline == 1 {  // 1st dataline, some special ?
				}
			}
	}

	// After sheet has written, setup style for cols

	if (exampleRow == nil ) { // if not used Excel ExampleRow, set col styles 
		useConfigCols(sheet, sheetMaxCols ,  cols  )
		SetColsDefaultStyle(sheet, cols)
	}
	if  myParam.debug>0  { fmt.Println("last line, setup cols:",sheetMaxCols,colnames) }

	if myParam.verbose>0 {
		fmt.Println()
		fmt.Println()
		}

	return sheetMaxCols, nil
}

// check if colnames include typing [d|i|n], default is text
func setColsType(colnames []string , coltypes *[]string ) {

	for colnr, colname := range colnames {
		myTypeStr := ""  // don't set text ..., it's default
		re := regexp.MustCompile(`\[i\]$`)
		myNewVal := re.ReplaceAllString(colname,"")
		if myNewVal != colname {
			myTypeStr = "int"
			colname=myNewVal
		}
		re = regexp.MustCompile(`\[d\]$`)
		myNewVal = re.ReplaceAllString(colname,"")
		if myNewVal != colname {
			myTypeStr = "date"
			colname=myNewVal
		}
		re = regexp.MustCompile(`\[n\]$`)
		myNewVal = re.ReplaceAllString(colname,"")
		if myNewVal != colname {
			myTypeStr = "float"
			colname=myNewVal
		}
		colnames[colnr]=colname
		*coltypes=append(*coltypes,myTypeStr)
	}
}

// write Row
func writeRowToXls(sheet *xlsx.Sheet, record []string, exampleRow *xlsx.Row, rownr int, colnames []string, coltypes []string, numOfCols int, linenr int  ) {

	var row *xlsx.Row
	var cell *xlsx.Cell
	//var celltype *xlsx.CellType

	var cellsLen int = 0
	if exampleRow != nil {
		cellsLen = len(exampleRow.Cells)
	}

	row = sheet.AddRow()
	//row.WriteSlice( &record , -1)


	if myParam.debug>1 { fmt.Println("Koko:",len(record)) }
	// write cells
	for colnr, colvalue := range record {
		if  colnr>=numOfCols { break }
		cell = row.AddCell()

		if rownr==1 { // init colwidth
			colsWidth[colnr]= float64(len(colvalue))
			}

		if exampleRow != nil && cellsLen > colnr {  // example row, use it
			if myParam.debug>0 { fmt.Println("ex") }
			writeCell(cell, exampleRow, colnr, colvalue, rownr, linenr )
		} else {  // no example row, so try to quess coltype, setCell do it
				if []rune(colvalue)[0] == '='  {   // input data include =formula syntax
					cell.SetFormula(colvalue[1:])
				} else {
					setCell(sheet,cell, colvalue, colnr, colnames[colnr], coltypes[colnr], linenr)
				}

				// length of value, set col width
				collen := float64(len(colvalue))
				if collen>colsWidth[colnr] { colsWidth[colnr]=collen }
		}
	} //for
}


// set cell value, return somete int value which tell some type of number or string
func setCellValue(cell *xlsx.Cell, v string , mType string) int {
	if mType == "" || mType == "int" { // try to set automatic Int using input string
		intVal, err := strconv.Atoi(v)
		if err == nil { // it's Int
			if intVal < 100000000000 { // Long numbers are displayed incorrectly in Excel
				cell.SetInt(intVal)
				return 1
			}
			cell.Value = v
			return 2
		}
	}

	if mType == "" || mType == "float" { // try to set automatic Float using input string
		floatVal, err := strconv.ParseFloat(v, 64)
		if err == nil { // It's float
			cell.SetFloat(floatVal)
			return 3
		}
	}

	// some string
	cell.Value = v
	return 99
}

// setCell set data in correct format.
func setCell(sheet *xlsx.Sheet, cell *xlsx.Cell, value string,  colNr int, colname string, coltype string, linecnt int) {


	var myType int

	var colPar *Col = myParam.confjson.Colskey[colname]

	if myParam.debug>0 {
		fmt.Println("----------------colPar:")
		spew.Dump(colPar)
		fmt.Println("-----------------------")
	}

	if colPar != nil {
			if linecnt == 1 && myParam.headerlines>0 { // headerline 
				coltype="text"  // headerline always text
			}
			if coltype == "" && colPar.Fldtype != "" { coltype=colPar.Fldtype }
	}


	// set value 1st
	myType=setCellValue(cell,value,coltype )

	// then setup style

	var cStyle *xlsx.Style
	var cNumFmt string
	var cValue string
	var intVal int
	var floatVal float64
	var err error
	var defStyle = xlsx.NewStyle()
	var defFont *xlsx.Font

	cStyle = cell.GetStyle()
	cNumFmt= cell.NumFmt
	cValue= cell.Value

	if colPar == nil   { // use default type formats
		switch coltype {
			case "int":	cell.SetFormat("0")
			case "date":	cell.SetFormat("d\\.m\\.yyyy;@")
			case "float":	cell.SetFormat("#,##0.00")  // 2 decimals
		}
		// default style
		defFont = xlsx.NewFont(myParam.fontsize, myParam.font)
		defStyle = xlsx.NewStyle()
		defStyle.Font = *defFont
		cell.SetStyle( defStyle )
	}


	// default has done if not used json config
	if myParam.defstyle == nil { return }

	*defStyle = *myParam.defstyle
	if myParam.debug>0 {
		fmt.Println(" Default Style")
		spew.Dump(myParam.defstyle)
		fmt.Println("")
		fmt.Println(" Col:(",colNr,") ",colname)
		fmt.Println(" Type:(",colNr,") ",coltype)
		fmt.Println("  Cell Style: %# v", pretty.Formatter(cStyle))
		fmt.Println("  Cell Value:", cValue, " Format:",cNumFmt)
		fmt.Println("  MyType Value:", myType)
		fmt.Println("  defStyle:", defStyle)
		}
	cell.SetStyle( defStyle )
	cStyle = cell.GetStyle()
	if myParam.debug>0 {
		fmt.Println("  After CelleStyle: %# v", pretty.Formatter(cStyle))
		}

	if colPar != nil   { // we have json config for this field
			var fontsize = myParam.fontsize
			if myParam.debug>0 { fmt.Println("    Check json config:",colname) }

			if linecnt == 1 && myParam.headerlines>0 { // headerline 
				coltype="text"  // headerline always text
			}
			if coltype == "" && colPar.Fldtype != "" { coltype=colPar.Fldtype }

			if myParam.debug>0 { fmt.Println("    Set type:",coltype) }
			switch coltype {
				case "text":
					cell.SetString(value)
				case "int":
					intVal, err = strconv.Atoi(value)
					if err == nil { // it's Int
						if intVal < 100000000000 { // Long numbers are displayed incorrectly in Excel
							if myParam.debug>0 { fmt.Println("  - Int:",value, intVal) }
							cell.SetInt(intVal)
						} else {
							if myParam.debug>0 { fmt.Println("  - Int iso:",value) }
							cell.SetValue(value)
						}
					} else {
						//cell.Value = v
						if myParam.debug>0 { fmt.Println("  - Int str:",value) }
						cell.SetValue(value)
					}
				case "float":
					floatVal, err = strconv.ParseFloat(value, 64)
					if err == nil { // It's float
						if myParam.debug>0 { fmt.Println("  - Float:",value, floatVal) }
						cell.SetFloat(floatVal)
					} else {
						if myParam.debug>0 { fmt.Println("  - Float Str:",value) }
						cell.SetValue(value)
					}
				case "date":
					cell.SetString(value)
					cell.SetFormat("d\\.m\\.yyyy;@")
			}  // set  type

			if colPar.Format != "" {
				if myParam.debug>0 { fmt.Println("    Col-",colname," format current:",cell.GetNumberFormat," new format:",colPar.Format) }
				cell.SetFormat(colPar.Format)
			}

			if colPar.Fontsize != "" {
				if myParam.debug>0 { fmt.Println("    Col-",colname," fontsize current:",cStyle.Font.Size," new fontsize:",colPar.Fontsize) }
				// overwrite fontsize default
				fontsize,_ = strconv.Atoi(colPar.Fontsize)
				cStyle.Font.Size=fontsize
			}

			if colPar.Width != "" {
				floatVal, err = strconv.ParseFloat(colPar.Width, 64)
				if err == nil {
					colsWidth[colNr]= (floatVal-2)
					sheet.SetColWidth(colNr+1,colNr+1,floatVal)
				}
			}

			if colPar.Font != "" {
				if myParam.debug>0 { fmt.Println("    Col-",colname," font current:",cStyle.Font.Name," new font:",colPar.Font ) }
				cStyle.Font.Name=colPar.Font
			}

			if colPar.Align != "" {
				if myParam.debug>0 { fmt.Println("    Col-",colname," align current:",cStyle.Alignment.Horizontal," new align:",colPar.Align) }
				cStyle.Alignment.Horizontal=colPar.Align
			}

			if colPar.Bold != "" {
				if myParam.debug>0 { fmt.Println("    Col-",colname," bold current:",cStyle.Font.Bold," new bold:",colPar.Bold) }
				if colPar.Bold == "true" { cStyle.Font.Bold=true }
				if colPar.Bold == "falsde" { cStyle.Font.Bold=false }
			}

			if colPar.Underline != "" {
				if myParam.debug>0 { fmt.Println("    Col-",colname," Underline current:",cStyle.Font.Underline," new Underline:",colPar.Underline) }
				if colPar.Underline == "true" { cStyle.Font.Underline=true }
				if colPar.Underline == "falsde" { cStyle.Font.Underline=false }
			}

			if colPar.Italic != "" {
				if myParam.debug>0 { fmt.Println("    Col-",colname," Italic current:",cStyle.Font.Underline," new Italic:",colPar.Italic) }
				if colPar.Italic == "true" { cStyle.Font.Italic=true }
				if colPar.Italic == "falsde" { cStyle.Font.Italic=false }
			}

			cell.SetStyle(cStyle)
			cStyle = cell.GetStyle()
			if myParam.debug>0 {
				fmt.Println("    After JSON CelleStyle: %# v", pretty.Formatter(cStyle))
			}
	}
}

// Write Cell using ExampleRow
func writeCell(cell *xlsx.Cell, exampleRow *xlsx.Row, colNr int, colString string, rownr int, datarow int) {

	var formula string
	var celltype xlsx.CellType

	formula = ""
	colStr := colString
	// get formula , it's empty string if no formula
	excell := exampleRow.Cells[colNr]
        formula = excell.Formula()
	celltype = excell.Type()
	cStyle := excell.GetStyle()
	cFormat := excell.GetNumberFormat()

	if ( excell.IsTime() ) {
		//cTime := excell.GetTime()
		}

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
	if myParam.debug>0 && datarow==1 { //oli 1  // printout some debug when we have 1st line
			fmt.Println("----------------excell:",colNr,colString)
			spew.Dump(exampleRow.Cells[colNr])
			fmt.Println("--------------------formula:")
			spew.Dump(formula)
			fmt.Println("--------------------cellType:")
			spew.Dump(celltype)
			fmt.Println("--------------------cStyle:")
			spew.Dump(cStyle)
			fmt.Println("--------------------cFormat:")
			spew.Dump(cFormat)
			fmt.Println("---------------------------")
	}

	// - allways 2 = xlsx.CellTypeNumeric ????
	switch celltype {
		// currently it looks that is always type 2 = CellTypeNumeric
		case xlsx.CellTypeNumeric:
			if myParam.debug>0 && datarow==1  {  fmt.Println("CellTypeNumeric",colNr)   }
			floatVal, err := strconv.ParseFloat(colStr, 64)
			if err == nil { // Float 
				if myParam.debug>0 && datarow==1 { fmt.Println(" -- float",floatVal,cFormat) }
				setCellValue(cell, colStr, "float")
				cell.SetFloatWithFormat(floatVal, cFormat)
			} else { // all other ...
				if myParam.debug>0 && datarow==1 { fmt.Println(" -- other Format ",cFormat) }
				setCellValue(cell, colStr,"")
				cell.SetFormat(cFormat)
			}
		default:
			if myParam.debug>0 && datarow==1  {  fmt.Println("celltype Default",colNr)   }
			setCellValue(cell, colStr, "")
		}

	cell.SetStyle(cStyle)

	celltype =  exampleRow.Cells[colNr].Type()
	cellwidth :=  exampleRow.Cells[colNr].Type()

	// if config file has set celltype, use it
	// ....
	if myParam.debug>0 && datarow==1   {
		fmt.Println("After set, Type:",celltype)
		fmt.Println("After set, Width:",cellwidth)
		}

}

// getCsvData read's data from CSV file.
func getCsvData(dataFileName string) (*csv.Reader, error) {

	dataFile, err := os.Open(dataFileName)
	if err != nil {
		return nil, cli.NewExitError("Problem with reading data from "+dataFileName, 11)
	}

	return csv.NewReader(dataFile), nil
}


//func useConfig(sheet *xlsx.Sheet, json string ) {
func useConfig(sheet *xlsx.Sheet ) {

	// reset maybe config has reset ... 
	xlsx.SetDefaultFont(myParam.fontsize,myParam.font)
	if myParam.debug>0 { fmt.Println("Set DefaultFont:",myParam.fontsize,myParam.font) }

}

// setup Cols style using config file values if exists, not using template xlsx sheet
func useConfigCols(sheet *xlsx.Sheet, numberOfCols int , colnames []string  ) {

	var col *xlsx.Col

        if myParam.debug>0  { // some debugs from cols
		fmt.Println("_______________________________________________________________________________________________________")
		fmt.Println("useConfigCols")
		fmt.Println("Cols debug")
		fmt.Println("Cols nr:",numberOfCols)
		fmt.Println("")
	}


	for colnr, _ := range colnames {
		sheet.SetColWidth(colnr+1,colnr+1,colsWidth[colnr]+2)  // default width for every Cols
		if myParam.debug>0 {
			fmt.Println("col width has set ",colnr)
			col = sheet.Cols.FindColByIndex(colnr+1)
			if col != nil {
				fmt.Println("col index work",colnr)
			} else {
				fmt.Println("col index not work",colnr)
			}
		}
	}

        if myParam.debug>0  { // some debugs from cols
		fmt.Println("useConfigCols END")
		fmt.Println("_______________________________________________________________________________________________________")
	}

	return
}

// setup cols default style
// currently not work with v2, need to setup cell allways ...
func SetColsDefaultStyle(sheet *xlsx.Sheet, values []string) {

	var col *xlsx.Col
	var defStyle *xlsx.Style
	var cStyle *xlsx.Style
	var defFont *xlsx.Font

	//defFont = &xlsx.Font{Size: 12, Name: "Verdana"}
	defFont = xlsx.NewFont(myParam.fontsize, myParam.font)
	defStyle = xlsx.NewStyle()
	defStyle.Font = *defFont

	if myParam.debug>0 {
		fmt.Println("")
		fmt.Println("SetColsDefaultStyle  ")
		fmt.Println("  Style: %# v", pretty.Formatter(defStyle))
		fmt.Println("  Sheet: %# v", pretty.Formatter(sheet))
		fmt.Println("  Spew sheet:")  // So good, real dump
		spew.Dump(sheet)
		fmt.Println("  ------")
		fmt.Println("  Cols:")
		fmt.Println("  Col: %# v", pretty.Formatter(sheet.Cols))
	}

	for colnr, _ := range values {
		if myParam.debug>0 { fmt.Println(" **** col id:",colnr) }
		col = sheet.Cols.FindColByIndex(colnr+1)
		if col == nil { break }
		if myParam.debug>0 { fmt.Println("  Col: %# v", pretty.Formatter(col)) }
		cStyle = col.GetStyle()
		if myParam.debug>0 && cStyle != nil {
			fmt.Println("  Style: %# v", pretty.Formatter(cStyle))
		} else {
			if myParam.debug>0 { fmt.Println("  no style") }
		}
		col.SetStyle(defStyle)
	}

	if myParam.debug>0 {
		fmt.Println("")
		fmt.Println("SetColsDefaultStyle done")
		fmt.Println("")
	}
}



