package main

import (
	"os"
	"fmt"
	"path/filepath"
	"encoding/json"
	_ "io"
"io/ioutil"
	"strconv"

	"github.com/urfave/cli"
	"github.com/tealeg/xlsx"
	_ "github.com/tidwall/gjson"
	"github.com/davecgh/go-spew/spew"    // nice dump system, help lot of understand structures, pointers, maps, ...
)

func initCommandLine(args []string) error {
	cli.NewApp()

	//cli.OsExiter = func(c int) {
	//	fmt.Fprintf(cli.ErrWriter, "error number %d\n", c)
	//}

	app := cli.NewApp()
	app.Name = "csv2xlsx"
	app.Usage = "Convert CSV data to XLSX - especially the big one. \n\n" +
		"Example: \n" +
		"   csv2xlsx --template example/template.xlsx --footer example/templatefooter.xlsx --sheet Sheet_1 --sheet Sheet_2 --row 2 --output result.xlsx data.csv data2.csv \n" +
		"   csv2xlsx.exe -t example\\template.xlsx -s Sheet_1 -s Sheet_2 -r 2 -o result.xlsx data.csv data2.csv \n"  +
		"   csv2xlsx.exe -sheetdefaultname MySheet -o result.xlsx data.csv  \n"  +
		"   csv2xlsx -d 0 -c ';' -t example/template5.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3.xlsx  example/data3.csv"

	app.Version = "0.2.2"
	app.ArgsUsage = "[file of file's list with csv data]"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "sheets, s",
			Usage: "sheet `names` in the same order like csv files. If sheet with that name exists, data is inserted to this sheet. Usage: -s AA -s BB ",
		},
		cli.StringFlag{
			Name:  "template, t",
			Value: "",
			Usage: "`path` to xlsx file with template output",
		},
		cli.StringFlag{
			Name:  "footer",
			Value: "",
			Usage: "`path` to xlsx file append to end of output",
		},
		cli.IntFlag{
			Name:  "row, r",
			Value: 0,
			Usage: "row `number` to use for create rows format. When '0' - not used. This row will be removed from xlsx file.",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "./output.xlsx",
			Usage: "path to result `xlsx file`",
		},
		cli.StringFlag{
			Name:  "config",
			Value: "",
			Usage: "`path` to config file, json format",
		},
		cli.StringFlag{
			Name:  "sheetdefaultname",
			Value: "Sheet",
			Usage: "Sheet default name, default is Sheet",
		},
		cli.StringFlag{
			Name:  "colsep, c",
			Value: ";",
			Usage: "column separator (default ';').",
		},
		cli.IntFlag{
			Name:  "debug, d",
			Value: 0,
			Usage: "debug 0 | 1.",
		},
		cli.IntFlag{
			Name:  "startrow",
			Value: 1,
			Usage: "start reading row `number` from csv, default 1st line = 1.",
		},
		cli.IntFlag{
			Name:  "headerlines",
			Value: 1,
			Usage: "`number` of headerlines in CSV, default 1.",
		},
		cli.IntFlag{
			Name:  "writeheaderlines",
			Value: 1,
			Usage: "writeheaderlines 0 | 1, default 1.",
		},
		cli.IntFlag{
			Name:  "verbose",
			Value: 0,
			Usage: "verbose, show rowcount 0 | 1, default 0.",
		},
		cli.StringFlag{
			Name:  "font, f",
			Value: "helvetica",
			Usage: "fonfamily name (default helvetica).",
		},
		cli.IntFlag{
			Name:  "fontsize",
			Value: 10,
			Usage: "default 10.",
		},
		cli.StringFlag{
			Name:  "formatdate",
			Value: "d.m.yyyy",
			Usage: "date format (default d.m.yyyy).",
		},
		cli.StringFlag{
			Name:  "formatnumber",
			Value: "#,##0.00",
			Usage: "number format (default #,##0.00).",
		},
	}

	//cli.Command{}

	app.Action = func(c *cli.Context) error {

		params, err := checkAndReturnParams(c)
		if err != nil {
			return err
		}
		myParam = *params // save param to the Global param struct
		return buildXls(c, params)
	}

	return app.Run(args)
}

func checkAndReturnParams(c *cli.Context) (*params, error) {
	p := &params{}

	output := c.String("output")
	if output == "" {
		return nil, cli.NewExitError("Path to output file not defined", 1)
	}

	output, err := filepath.Abs(output)
	if err != nil {
		return nil, cli.NewExitError("Wrong path to output file", 2)
	}
	p.output = output

	//JUI
	p.colsep = c.String("colsep")
	if p.colsep == "" {
		return nil, cli.NewExitError("Column separator not defined", 1)
	}

	//

	p.input = make([]string, len(c.Args()))
	for i, f := range c.Args() {
		filename, err := filepath.Abs(f)
		if err != nil {
			return nil, cli.NewExitError("Wrong path to input file "+filename, 3)
		}
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return nil, cli.NewExitError("Input file does not exist ( "+filename+" )", 4)
		}

		p.input[i] = filename
	}

	//

	p.row = c.Int("row")
	p.debug = c.Int("debug")
	p.headerlines = c.Int("headerlines")
	p.startrow = c.Int("startrow")
	p.writeheaderlines = c.Int("writeheaderlines")
	p.verbose = c.Int("verbose")
	p.fontsize = c.Int("fontsize")
	p.font = c.String("font")
	p.formatdate = c.String("formatdate")
	p.formatfloat = c.String("formatnumber")
	p.sheetdefaultname = c.String("sheetdefaultname") + " %d"  // Sheet % d 

	p.sheets = c.StringSlice("sheets")
	configFile := c.String("config")

	//

	xlsxTemplate := c.String("template")
	xlsxFooter := c.String("footer")

	if xlsxTemplate != "" {
		xlsxTemplate, err = filepath.Abs(xlsxTemplate)
		if err != nil {
			return nil, cli.NewExitError("Wrong path to template file", 5)
		}
		if _, err := os.Stat(xlsxTemplate); os.IsNotExist(err) {
			return nil, cli.NewExitError("Template file does not exist ( "+xlsxTemplate+" )", 6)
		}
		p.xlsxTemplate = xlsxTemplate
	}

	if xlsxFooter != "" {
		xlsxFooter, err = filepath.Abs(xlsxFooter)
		if err != nil {
			return nil, cli.NewExitError("Wrong path to footer file", 5)
		}
		if _, err := os.Stat(xlsxFooter); os.IsNotExist(err) {
			return nil, cli.NewExitError("Footer file does not exist ( "+xlsxFooter+" )", 6)
		}
		p.xlsxFooter = xlsxFooter
	}


	if p.row != 0 && xlsxTemplate == "" {
		return nil, cli.NewExitError("Defined `row template` without xlsx template file", 7)
	}

	p.json="" // default, no json readed from config file

	if configFile != "" {
		if (p.debug>0) { fmt.Println("conf",configFile) }
		configFile, err = filepath.Abs(configFile)
		if err != nil {
			return nil, cli.NewExitError("Wrong path to config file", 5)
		}
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			fmt.Println("fpath",configFile,"err")
			return nil, cli.NewExitError("Config file does not exist ( "+configFile+" )", 6)
		}
		p.config=configFile

		// Open our jsonFile - configuration
		jsonFile, err := os.Open(p.config) // overwrite commanlines arguments if exists
		// if we os.Open returns an error then handle it
		if err != nil {
				fmt.Println(err)
		p.config=""  // don't use ... file not opened
		}

		// handel json config if opened
		if p.config != "" && jsonFile != nil {
				// defer the closing of our jsonFile so that we can parse it later on
				defer jsonFile.Close()
				// read our opened xmlFile as a byte array.
				byteValue, _ := ioutil.ReadAll(jsonFile)
				p.bjson = byteValue  // 
				p.json = string(byteValue)  // gjson library like more strings as []byte

				/*
				if p.json != "" {
					value := gjson.Get(p.json, "font")
					if value.Index>0  { p.font = value.Str }
					value = gjson.Get(p.json, "fontsize")
					if value.Index>0  { p.fontsize,_ = strconv.Atoi(value.Str) }
					value = gjson.Get(p.json, "verbose")
					if value.Index>0  { p.verbose,_ = strconv.Atoi(value.Str) }
					value = gjson.Get(p.json, "headerlines")
					if value.Index>0  { p.headerlines,_ = strconv.Atoi(value.Str) }
					value = gjson.Get(p.json, "writeheaderlines")
					if value.Index>0  { p.writeheaderlines,_ = strconv.Atoi(value.Str) }
				}
				*/

				// parse json to the structs ... nice
				err = json.Unmarshal(byteValue,&p.confjson)
				if err != nil {
					fmt.Println("json parser error:",err)
					os.Exit(5)
				}

				if p.confjson.Font != "" { p.font = p.confjson.Font  }
				if p.confjson.Fontsize != "" { p.fontsize,_ = strconv.Atoi(p.confjson.Fontsize)  }
				if p.confjson.Verbose != "" { p.verbose,_ = strconv.Atoi(p.confjson.Verbose)  }
				if p.confjson.Headerlines != "" { p.headerlines,_ = strconv.Atoi(p.confjson.Headerlines)  }
				if p.confjson.Writeheaderlines != "" { p.writeheaderlines,_ = strconv.Atoi(p.confjson.Writeheaderlines)  }
				if p.confjson.FormatDate != "" { p.formatdate = p.confjson.FormatDate  }
				if p.confjson.FormatFloat != "" { p.formatfloat = p.confjson.FormatFloat  }


				// JSON cols array
				confCols := p.confjson.Cols
				// init map using indexing - colname, array value is pointer to the Col structure
				p.confjson.Colskey = make(map[string]*Col)
				if p.debug>0 {
					fmt.Println("DEBUG")
					spew.Dump(p.confjson)
				}

				// indexing cols using colname
				// array/map of cols, need ptr to the Col struct element
				var key string
				var pcol *Col

				// loop cols array
				for id,colvalue := range confCols {
					if p.debug > 0 { fmt.Println("  col:",id,colvalue.Colname,colvalue) }
					// pointer to the new Col object
					pcol=new(Col)
					// take cols array one "line" and copy it to the some mem addr
					// of course we could use also original cols array, but here we copy col data to the new object
					*pcol=colvalue  // copy col structure to new object (pointer)
					if p.debug > 0 {
						fmt.Println("_added key_______________________")
						fmt.Println("col",id)
						fmt.Println("name",pcol.Colname)
						spew.Dump(pcol)
					}
					// colname is index, cols array value is pointer to the Col-struct element
					key=pcol.Colname
					if p.debug > 0 { fmt.Println("name",pcol.Colname) }
					// add to the map key + pointer to the object, pcol is pointer
					p.confjson.Colskey[key] = pcol
				}

				// now cols include all col nodes
				// check result
				if p.debug > 0 {
					fmt.Println("==========================Cols values:")
					//spew.Dump(p.confjson.Colskey)
					for key, col := range p.confjson.Colskey {
						fmt.Println("col:",key,col)
					}
					fmt.Println("======================================")
				} // json config parser
		}

		// setup defaults
		p.deffont = xlsx.NewFont(p.fontsize, p.font)
		p.defstyle = xlsx.NewStyle()
		p.defstyle.Font = *p.deffont
		if p.debug>0 {
			fmt.Println("DefFont:",p.deffont)
			fmt.Println("DefStyle:",p.defstyle)
		}

	} // configfile

	return p, nil
}

type params struct {
	output string
	config string
	input  []string
	xlsxTemplate string
	xlsxFooter string
	sheets []string
	row    int
	colsep string
	headerlines int
	startrow int
	writeheaderlines int
	verbose int
	fontsize int
	font string
	formatdate string
	formatfloat string
	json string
	bjson []byte
	sheetdefaultname string
	confjson Config  // parsed json
	debug int
	deffont *xlsx.Font
	defstyle *xlsx.Style
}



// Global param
var myParam  params

type Config struct {
        Fontsize		string	`json:"fontsize"`
        Font			string  `json:"font"`
        FormatDate		string  `json:"formatdate"`
        FormatFloat		string  `json:"formatnumber"`
        Verbose			string  `json:"verbose"`
        Headerlines		string  `json:"headerlines"`
        Writeheaderlines	string  `json:"writeheaderlines"`
        Cols []Col			`json:"cols"`
        Colskey map[string]*Col  // - col indexing using colname
        }

type Col struct {
        Colname		string `json:"colname"`
        Id		int `json:"id"`
        Fldtype		string `json:"type"`
        Font		string `json:"font"`
        Format		string `json:"format"`
        Width		string `json:"width"`
        Fontsize	string `json:"fontsize"`
        Bold		string `json:"bold"`
        Underline	string `json:"underline"`
        Italic		string `json:"italic"`
        Align		string `json:"align"`
}

