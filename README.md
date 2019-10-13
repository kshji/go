# Go Language #

Here is some tools programmed using Go Language.

I'm not used so much Go, but when I like to make some nice ex. commandtool to *nix and Windows, then my selection is Go.
Nice cross compiling, one binary file to install. No library hasling in the destination. 

I am in the process of learning Go and therefore
I am sure there are much better, more Go-idiomatic ways to achieve this functionality. If you have feedback on how to improve
the code or want to contribute, please do not hesitate to do so. I'd really like to improve my GO skills and learn things.

Go is more hoppy and nice tool to make compact program to install *nix+Windows.

In real life I need lot of convert CSV output to the Excel and
read input data from Excel. CSV/Json/XML formats are my working formats, because ex. Postgresql like to use COPY command to import/export.
Also awk-command,ksh,... like to use CSV files.

  * [xlsx2csv] (https://github.com/tealeg/xlsx2csv) , I have only added Sheetname support to export CSV. Very simple and fast Excel to Csv converter.
  * [xlsx] (https://github.com/tealeg/xlsx), nice XLSX library to read and write Excel files.
  * [xlsx doc] (https://godoc.org/github.com/tealeg/xlsx), xlsx documentation
  * [csv2xlsx] (https://gitlab.com/DerLinkshaender/csv2xlsx) This converter has given some ideas. Thanks. 
  * [csv2xlsx] (https://github.com/mentax/csv2xlsxts)  I have took this source and updated it. 


## My Repo ##
   * [Awk] (https://github.com/kshji/awk)
   * [Ksh] (https://github.com/kshji/ksh)


## HELP
  Actual version always on  csv2xlsx -h or csv2xlsx help

### NAME:
   csv2xlsx - Convert CSV data to xlsx - especially the big one.

### Speed:

   csv with 50k rows, 5 MB, with xlsx template - 5s


   (On Windows 10, WSL Ubuntu 18.04 )

### Example:

```bash
csv2xlsx --template example/template.xlsx --sheet Sheet_1 --sheet Sheet_2 --row 2 --output result.xlsx data.csv data2.csv
csv2xlsx.exe -t example/template.xlsx -s Sheet_1 -s Sheet_2 -r 2 -o result.xlsx data.csv data2.csv
```

### USAGE:

    csv2xlsx [global options] command [command options] [file of file's list with csv data]

#### VERSION:
   0.2.1

#### GLOBAL OPTIONS:

```
--sheets names, -s names          sheet names in the same order like csv files. If sheet with that name exists, data is inserted to this sheet. Usage: -s AA -s BB
--template path, -t path          path to xlsx file with template output
--row number, -r number           template row number to use for create rows format. When '0' - not used. This row will be removed from xlsx file. (default: 0)
--output xlsx file, -o xlsx file  path to result xlsx file (default: "./output.xlsx")
--headerlines number		  how many headerlines in CSV, default 1
--writeheaderlines 0|1	  	  write headerlines to the Excel, default 1, yes. If templates include headers, then set this 0.
--startrow number		  Default is start csv reading from line 1. If not like import headerline, then set this ex. 2
--help, -h                        show help
--debug 0|1, -h                   debug level 0 | 1, default 0.
--verbose 0|1                     default 0. Show rownumber when processing csv files.
--version, -v                     print the version
```   


## Download

Original version:
Download from [releases section on GitHub](https://github.com/mentax/csv2xlsx/releases)   

My updated version
Download from [releases section on GitHub](https://github.com/kshji/go/csv2xlsx)   

