
package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"unicode/utf8"
	"io"
	"strconv"
	"regexp"
	"strings"
	"time"
	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
	"github.com/davecgh/go-spew/spew"
	"github.com/kr/pretty"
)

// SheetNamesTemplate define name's for new created sheets.
//var SheetNamesTemplate = "Sheet %d"

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

	dbg("Set DefaultFont:",myParam.fontsize,myParam.font)

	var sheetFooter *xlsx.Sheet = nil
	var xlRow *xlsx.Row
	var exampleRow *xlsx.Row
	// parse all env variables using indexed array - map - hash table
	var envvars = make(map[string]string)
	var varstr string = ""

	// loop all environment variables and add to the map - hash table - maybe need to expand if template include {XX} 
	for _, e := range os.Environ() {
                pair := strings.SplitN(e, "=", 2)
                dbg(e," - ",pair[0],"=",pair[1])
                varstr = "{"+pair[0]+"}"  // add {  } = string comparing is easier, because template using {}
                envvars[varstr] = pair[1]
	}

	regexprule := regexp.MustCompile(`{([^}]+)}`)   // {XXX}
	if myParam.debug > 0 {
		fmt.Println("RegExpRule:")
		spew.Dump(regexprule)
	}



	for i, dataFileName := range dataFiles {

		sheet, err := getSheet(xlFile, sheetNames, i)
		if err != nil {
			// make copy from template sheet 0 // waiting my freetime ....
			//sheet, err := getSheet(xlFile, sheetNames, i)
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

			sheet.Rows = append(sheet.Rows[:exampleRowNumber-1], sheet.Rows[exampleRowNumber:]...)
			//parse template, expand variables if used
			//parseSheet(sheet, exampleRowNumber-1,regexprule ,envvars)
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
		parseSheet(sheet, (-1) ,regexprule ,envvars)

	}

	return nil
}

func getSheet(xlFile *xlsx.File, sheetNames []string, i int) (sheet *xlsx.Sheet, err error) {

	var sheetName string
	if len(sheetNames) > i {
		sheetName = sheetNames[i]
	} else {
		//sheetName = fmt.Sprintf(SheetNamesTemplate, i+1)
		sheetName = fmt.Sprintf(myParam.sheetdefaultname, i+1)
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


// loop Sheet and expand variables if exists, if maxRow < 0 = all rows
func parseSheet(sheet *xlsx.Sheet, maxRow int,re *regexp.Regexp, variables  map[string]string ) {

	var newvalue string
	for rownr, row := range sheet.Rows {
		if myParam.debug>0 { fmt.Println("examplerow:",rownr,row) }
		if rownr >= maxRow && maxRow>0 { return }
		//var vals []string
		if row != nil {
			for cellnr, cell := range row.Cells {
				value, err := cell.FormattedValue()
				if err != nil {
					//vals = append(vals, err.Error())
				}
				newvalue = Expand(value,re, variables)
				if myParam.debug>0 { fmt.Println("  ex:",rownr,cellnr,value,newvalue) }
				if newvalue != value { cell.Value=newvalue }
			}
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
		fmt.Println("formatnumber",myParam.formatfloat)
		fmt.Println("formatdate",myParam.formatdate)
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
	var formulastr string
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
				if colvalue != "" && []rune(colvalue)[0] == '='  {   // input data include =formula syntax
					formulastr = colvalue[1:]
					//cell.SetFormula(colvalue[1:])
					// maybe need to check formulas which include " chars ...
					cell.SetFormula(formulastr)
				} else {
					setCell(sheet,cell, colvalue, colnr, colnames[colnr], coltypes[colnr], linenr)
				}

				// length of value, set col width
				collen := float64(len(colvalue))
				if collen>colsWidth[colnr] { colsWidth[colnr]=collen }
		}
	} //for
}


// set cell value, return some int value which tell some type of number or string
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
			// set format
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
	var date time.Time

	cStyle = cell.GetStyle()
	cNumFmt= cell.NumFmt
	cValue= cell.Value

	if colPar == nil   { // use default type formats
		if myParam.debug>0 { fmt.Println(" use default types:", myType, coltype ) }
		switch myType { // automatic type of value
			case 3:		cell.SetFormat(myParam.formatfloat)
		}
		switch coltype {
			case "int":	cell.SetFormat("0")
			case "date":
					//cell.SetFormat(myParam.formatdate) // cell.SetFormat("d\\.m\\.yyyy;@")
					//cell.SetFormat(myParam.formatdate) // cell.SetFormat("d.m.yyyy")
					date, err = time.Parse("2006-01-02",value)
					if err == nil { // it's time
						dateFormat := myParam.formatdate
						dateFormat=ConvertTimeFormat(dateFormat)
						value=date.Format(dateFormat)
						setCellValue(cell, value,"")
						//cell.SetFormat(myParam.formatdate)
						if myParam.debug>0 { fmt.Println("Date:",date,value,dateFormat) }
					} else {
						setCellValue(cell, value,"")
						cell.SetFormat(myParam.formatdate)
						if myParam.debug>0 { fmt.Println("Ei Date:",value,myParam.formatdate) }
						//cell.Value=FormatTime(date)
					}
			case "float":	cell.SetFormat(myParam.formatfloat) // cell.SetFormat("#,##0.00")  // 2 decimals
		}
		// default style
		defFont = xlsx.NewFont(myParam.fontsize, myParam.font)
		defStyle = xlsx.NewStyle()
		defStyle.Font = *defFont
		cell.SetStyle( defStyle )
	}


	// default has done if not used json config
	if myParam.defstyle == nil { return }

	// use config setup
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
					//cell.SetString(value)
					//cell.SetFormat("d\\.m\\.yyyy;@")
					//cell.SetFormat("d.m.yyyy")
					date, err = time.Parse("2006-01-02",value)
					if err == nil { // it's time
						dateFormat := myParam.formatdate
						if colPar.Format != "" { dateFormat=colPar.Format }
						dateFormat=ConvertTimeFormat(dateFormat)
						value=date.Format(dateFormat)
						setCellValue(cell, value,"")
						if myParam.debug>0 { fmt.Println("parDate:",date,value,dateFormat) }
					} else {
						setCellValue(cell, value,"")
						if myParam.debug>0 { fmt.Println("parNoDate:",value,myParam.formatdate) }
					}

			}  // set  type

			if colPar.Format != "" && coltype != "date" {
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
	var cTime time.Time
	var cellIsTime bool = false
	var date time.Time
	var err error

	formula = ""
	colStr := colString
	// get formula , it's empty string if no formula
	excell := exampleRow.Cells[colNr]
        formula = excell.Formula()
	celltype = excell.Type()
	cStyle := excell.GetStyle()
	cFormat := excell.GetNumberFormat()
	if cellIsTime { celltype = xlsx.CellTypeDate }


	if ( excell.IsTime() ) {
		cTime, _ = excell.GetTime(true)  // date1904
		if myParam.debug>0 { fmt.Println(" - Timetype",cTime) }
		// template cell is time but value is little difficult, because Excel turn to the ISO format ...
		// can't found format setup as it's in Excel ...
		// SetDate(t time.Time)
		cellIsTime=true
		// testing how it works ...
		celltype=xlsx.CellTypeDate
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
		case xlsx.CellTypeDate:
			if myParam.debug>0 && datarow==1  {  fmt.Println("CellTypeDate",colNr)   }
			date, err = time.Parse("2006-01-02",colStr)
			if myParam.debug>0 { fmt.Println(" -- time :",date) }
			//colStr=date.Format("2006-01-02")
			if err == nil { // it's time  - this is not so easy case ...
				//colStr=date.Format(cFormat)
				//dateFormat=ConvertTimeFormat(cFormat)
				//cell.SetDate( date)
				dateFormat := myParam.formatdate
				dateFormat=ConvertTimeFormat(dateFormat)
				colStr=date.Format(dateFormat)
				if myParam.debug>0 && datarow==1 { fmt.Println(" -- date Format ",myParam.formatdate,dateFormat,cFormat) }
				setCellValue(cell, colStr,"")
				cell.SetFormat(cFormat)
			} else {
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
	//cellwidth :=  exampleRow.Cells[colNr].Width

	// if config file has set celltype, use it
	// ....
	if myParam.debug>0 && datarow==1   {
		fmt.Println("After set, Type:",celltype)
		//fmt.Println("After set, Width:",cellwidth)
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

// convert printf +timeformats to the go time formats
func ConvertTimeFormat(src string) string {
	result := src
	result=strings.Replace(result, "mm", "01", -1)
	result=strings.Replace(result, "m", "1", -1)
	result=strings.Replace(result, "dd", "02", -1)
	result=strings.Replace(result, "d", "2", -1)
	result=strings.Replace(result, "yyyy", "2006", -1)
	result=strings.Replace(result, "yy", "06", -1)
	result=strings.Replace(result, "\\", "", -1)
	//result=strings.Replace(result, ";@", "", -1)  // maybe need to split ; and use only 1st arg
	return result
}

// replace variables syntax {XXX} using XXX value
func Expand(instr string, regexprule *regexp.Regexp, variables  map[string]string )  string {
	// input include 0-n {variables} and variables are of course not allways the same 
	// loop all matches from input match by match and make unique replace all of those
	result:=regexprule.ReplaceAllStringFunc(instr,
		func(inputstr string) string {
			// search if there is {variable} ..., find 1st match
			match := regexprule.FindString(inputstr)
			if  match == "" { return inputstr }
			replace := variables[match]
			if  replace == "" { return inputstr }
			// replace it
			return regexprule.ReplaceAllString(inputstr, replace)
		})
	return result
}

func dbg(args ...interface{}) {

	if myParam.debug<1 { return }
        fmt.Printf("Debug ")
        for _,arg := range args {
                fmt.Printf("%v",arg)
        }
        fmt.Println()
}


func dbgDump(str interface{}, variable interface{}) {
	if myParam.debug>0 { return }
        fmt.Println("/*")
        fmt.Println("  dbgDump ",str)
        spew.Dump(variable)
        fmt.Println("*/")
}
