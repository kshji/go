

[![Build Status](https://travis-ci.org/mentax/csv2xlsx.svg?branch=master)](https://travis-ci.org/mentax/csv2xlsx)
[![GoDoc](https://godoc.org/github.com/mentax/csv2xlsx?status.svg)](https://godoc.org/github.com/mentax/csv2xlsx)
[![codebeat badge](https://codebeat.co/badges/1b57272c-e0fa-4a14-93b5-3586e192fdb3)](https://codebeat.co/projects/github-com-mentax-csv2xlsx-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mentax/csv2xlsx)](https://goreportcard.com/report/github.com/mentax/csv2xlsx)
[![Coverage](https://gocover.io/_badge/github.com/mentax/csv2xlsx)](http://gocover.io/github.com/mentax/csv2xlsx)

# csv 2 xlsx


## HELP
  Actual version always on  csv2xlsx -h or csv2xlsx help

### NAME:
   csv2xlsx - Convert CSV data to xlsx - especially the big one.

### Speed:

   csv with 50k rows, 5 MB, with xlsx template - 5s


   (On Windows 10, Intel i7, WSL Ubuntu 18.04 )


### Example:

```bash
csv2xlsx --template example/template.xlsx --sheet Sheet_1 --sheet Sheet_2 --row 2 --output result.xlsx data.csv data2.csv
csv2xlsx.exe -t example/template.xlsx -s Sheet_1 -s Sheet_2 -r 2 -o result.xlsx data.csv data2.csv
```

### USAGE:

    csv2xlsx [global options] command [command options] [file of file's list with csv data]

#### VERSION:
   0.2.2

#### GLOBAL OPTIONS:

```
--sheets names, -s names          sheet names in the same order like csv files. If sheet with that name exists, data is inserted to this sheet. Usage: -s AA -s BB
--template path, -t path          path to xlsx file with template output
--row number, -r number           template row number to use for create rows format. When '0' - not used. This row will be removed from xlsx file. (default: 0)
--output xlsx file, -o xlsx file  path to result xlsx file (default: "./output.xlsx")
--headerlines number              how many headerlines in CSV, default 1
--writeheaderlines 0|1            write headerlines to the Excel, default 1, yes. If templates include headers, then set this 0.
--startrow number                 Default is start csv reading from line 1. If not like import headerline, then set this ex. 2
--help, -h                        show help
--debug 0|1, -h                   debug level 0 | 1, default 0.
--verbose 0|1                     default 0. Show rownumber when processing csv files.
--version, -v                     print the version
```

## LICENSE
	mentax has done excellent packet. I have only add some extensions. Enjoy.
	[License] (https://github.com/mentax/csv2xlsx/blob/master/LICENSE)

## Download

Original version:
Download from [releases section on GitHub](https://github.com/mentax/csv2xlsx/releases)

My updated version
Download from [releases section on GitHub](https://github.com/kshji/go/csv2xlsx)

