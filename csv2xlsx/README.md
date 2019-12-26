
# csv2xlsx again
## csv to xlsx - csv to Excel
This **csv2xlsx** is little more as traditional convert csv using some format to the Excel.
 - you can use Excel template to tell format
 - you can use json config file to tell output rule like font, format, ...
 - csv can include formulas
 - template can include indirect addressing formulas
 - you can use environment variables in templates AND csv

## ORIGINAL VERSIONS

  * [mentax csv2xlsx](https://github.com/mentax/csv2xlsx) I have took this source and updated it.
  * [DerLinkshaender csv2xlsx](https://gitlab.com/DerLinkshaender/csv2xlsx) This converter has given some ideas. Thanks.
  * [tealeg xlsx] (https://github.com/tealeg/xlsx) I use tealeg Xlsx library
  * [Excelize] (https://github.com/360EntSecGroup-Skylar/excelize ) is interesting newer Xlsx library
  * [Plandem] (https://github.com/plandem/xlsx) An other Xlsx library for Go

## HELP
  Actual version always on  csv2xlsx -h or csv2xlsx help

### NAME:
   csv2xlsx - Convert CSV data to xlsx - especially the big one and/or using Excel templates to format output

### Speed:

   csv with 50k rows, 3.4 MB, with xlsx template with footer, every line include variable - 3.7 s

   (On Windows 10, Intel i7, WSL Ubuntu 18.04 )


### Example:

```bash
csv2xlsx --template example/template.xlsx --sheet Sheet_1 --sheet Sheet_2 --row 2 --output result.xlsx example/data.csv example/data2.csv
csv2xlsx.exe -t example/template.xlsx -s Sheet_1 -s Sheet_2 -r 2 -o result.xlsx example/data.csv example/data2.csv
csv2xlsx -d 0 -c ';' -t example/template5.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3.xlsx  example/data3.csv

# remove header + using template and footer template 
csv2xlsx -d 0 -c ';' -t template5.xlsx --footer template5footer.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3a.xlsx  data3.csv 
```

You can use also formulas in template or in csv. Csv formulas overwrite template formulas.
Look examples formula col.

#### Example include templates and screenshots

```bash
# use template template5.xlsx Sheet Sh2, and footer template template5footer.xlsx , row 5 is data example row
# input data.csv including headerline and not write it
# result to the file result.xlsx
csv2xlsx -c ';' -t template5.xlsx --footer template5footer.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o result.xlsx  data.csv
```
Result:
<img src="https://raw.githubusercontent.com/kshji/go/master/csv2xlsx/example/example.result.png?width=600&button=false" />
Template:
<img src="https://raw.githubusercontent.com/kshji/go/master/csv2xlsx/example/template_example.png?width=600&button=false" />
Csv:
<img src="https://raw.githubusercontent.com/kshji/go/master/csv2xlsx/example/data3.csv.png?width=600&button=false" />
Footer template:
<img src="https://raw.githubusercontent.com/kshji/go/master/csv2xlsx/example/template_footer_example.png?width=600&button=false" />

#### VERSION:
   2019-12-26

#### OPTIONS:

```
--sheets names, -s names          sheet names in the same order like csv files. If sheet with that name exists, data is inserted to this sheet. Usage: -s AA -s BB
--sheetdefaultname		  Sheet default name, default is Sheet (+ %d )
--template path, -t path          path to xlsx file with template output
--row number, -r number           template row number to use for create rows format. When '0' - not used. This row will be removed from xlsx file. (default: 0)
--footer footer_template_path     path to footer xlsx file - footer template 
--output xlsx file, -o xlsx file  path to result xlsx file (default: "./output.xlsx")
--colsep char, c char             column separator (default ';')
--headerlines number              how many headerlines in CSV, default 1
--writeheaderlines 0|1            write headerlines to the Excel, default 1, yes. If templates include headers, then set this 0.
--startrow number                 Default is start csv reading from line 1. If not like import headerline, then set this ex. 2
--config jsonconfigfile           config file, json format: default font, columns defination, used without templates
--formatnumber "#,##0.00"	  format of number cols, default "#,##0.00"
--formatdate ""d.m.yyyy"	  format of date cols, default "d.m.yyyy"
--help, -h                        show help
--debug 0|1, -d 0|1               debug level 0 | 1, default 0.
--verbose 0|1                     default 0. Show rownumber when processing csv files.
--version, -v                     print the version
```

#### CSV special
If headerline columnname ending using [d] or [i] or [n], then column typing has used, not default.
This need little development so that user can tell also default format. Currently format is builtin.
##### [d] date format yyyy-mm-dd
##### [i] integer 
##### [n] float format 0,00

##### Formulas

If csv cell start using symbol **=**, the cell will be formula, not value.

Example:
- =J:J+K:K  sum of column J and K in this line
- =K1*J:J   Cell K1 multiply value of cell J in this line

#### XLSX template and csv data special, expand environment variables 
Expand environment variables if exists in result sheet.

If cell include labeled string like {HOME} or {PATH} or any other environment variable name,
those will replace value of variable.

##### Formatting date
You can use Go supported format or also ex. yyyy-mm-dd, d.m.yyyy, ...
[Go time format](https://golang.org/src/time/format.go)

##### Formatting number
You can used xlsx library number formats, same as Excel use.

## TODO
 - 

## LICENSE

[mentax](https://github.com/mentax/) has done excellent packet. I have only add some extensions. Enjoy.

[License](https://github.com/mentax/csv2xlsx/blob/master/LICENSE)

[License](https://github.com/kshji/go/csv2xlsx/blob/master/LICENSE)

## Download

Original version:
Download from [releases section on GitHub](https://github.com/mentax/csv2xlsx/releases)

My updated version
Download from [releases section on GitHub](https://github.com/kshji/go/tree/master/csv2xlsx/build)

