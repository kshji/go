
# xlsx to csv
##  xlsx2csv 


## ORIGINAL VERSIONS

  * [tealeg xlsx2csv](https://github.com/tealeg/xlsx2csv) I have took this source and updated it.

## HELP
  Actual version always on  csv2xlsx -h or csv2xlsx help

### NAME:
   xlsx2csv - Convert XLSX files to csv 


### Example:

```bash
# read 1st sheet from my.xlsx
xlsx2csv -d ";" -f my.xlsx -i 0  > my.csv

#  read sheet named Sheetname, starting from row 3. 1st row is 1
xlsx2csv -d ";" -f my.xlsx -s "SheetName"  -r 3 > my.csv

```

You can use also formulas in template or in csv. Csv formulas overwrite template formulas.
Look examples formula col.

### USAGE:

    xlsx2csv -f xlsx_file [options] > file.csv

#### VERSION:

#### OPTIONS:

```
	-d CHAR		column delimiter char, default ;
	-f file.xlsx	input xlsx file
	-i number	sheet number to read, 1st is 0 and it's default
	-s sheetname	use sheetname, not sheet number
	-r startrow	default is to start from row 1

```


## LICENSE


[tealeg](https://github.com/tealeg/xlsx2csv/) has done excellent packet xlsx and also xsls2csv. I have add some extensions. Enjoy.

[License](https://raw.githubusercontent.com/tealeg/xlsx2csv/master/LICENSE)

## Download

Original version:
Download from [releases section on GitHub]((https://github.com/tealeg/xlsx2csv/)

My updated version
Download from [releases section on GitHub](https://github.com/kshji/go/tree/master/xlsx2csv/build)

