
# xlsxsheetcopy
##   xlsxsheetcopy


## HELP
  Actual version always on  xlsxsheetcopy -h or xlsxsheetcopy help

### NAME:
   xlsxsheetcopy - Copy sheet 


### Example:

```bash
# read 1st sheet from my.xlsx , make copy to the sheet Copied
# save result to the same file
# => xlsx can include lot of sheets, you only copy one sheet, it will be the last sheet
xlsxsheetcopy -f my.xlsx -i 0 -n Copied 

#  read sheet named csvsheet, make copy to the sheet Copied2
xlsx2csv -f my.xlsx -s csvsheet -n Copied2

```


### USAGE:

    xlsxsheetcopy -f update.xlsx -i 0 -n CopiedSheetName

#### VERSION:

#### OPTIONS:

```
	-f file.xlsx	xlsx file
	-i number	sheet number to read, 1st is 0 and it's default
	-s sheetname	use sheetname, not sheet number
	-n newsheetname name of new sheet, it'll be the last sheet in xlsx file

```


## LICENSE

[License](https://github.com/kshji/go/blob/master/xlsxsheetcopy/LICENSE)

## DOWNLOAD
Download from [releases section on GitHub](https://github.com/kshji/go/tree/master/xlsxsheetcopy/build)

