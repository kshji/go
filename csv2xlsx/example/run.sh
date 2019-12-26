#!/bin/sh


# Basic
../csv2xlsx --verbose 1 --font ubuntu --fontsize 15 -o data1.xlsx data1.csv
../csv2xlsx --font ubuntu --fontsize 16 -o data2.xlsx data2.csv
../csv2xlsx -o data1b.xlsx data1.csv
# name the sheets MySheet 1, MySheet 2, ...
../csv2xlsx --sheetdefaultname "MySheet" -o data1d.xlsx data1.csv
../csv2xlsx --formatdate "yyyy.mm.dd" --formatnumber "#,##0.0" -o data1c.xlsx data1.csv
# Use columnames typing
../csv2xlsx -o data4.xlsx data4.csv
../csv2xlsx --formatdate "yyyy.mm.dd" --formatnumber "#,##0.0" -o data4b.xlsx data4.csv

# Basic config
../csv2xlsx -d 0 --config data1.cfg --verbose 1 --font ubuntu --fontsize 15 -o data1d.xlsx data1.csv

# Multisheet import
echo "data2 multisheet, using template"
../csv2xlsx -d 0 -c ';' -t template.xlsx -r 4 -s Sh1 -s Sh2 -s NewSheet -o data2.multi.xlsx data2.csv data2.csv data2.csv

echo "data3 using jsonconfig"
../csv2xlsx -d 0 -c ';' --config data3.cfg  --headerlines 1 --writeheaderlines 1 -s Sh2 -o data3.xlsx  data3.csv

echo "data3 remove header, using template"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3a.xlsx  data3.csv

echo "data3 remove header, using template + footer"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --footer template5footer.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3a.xlsx  data3.csv 

echo "data3b not remove header, using template"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 1 -r 5 -s Sh2 -o data3b.xlsx  --verbose 1 data3.csv

echo "data3c remove header, use template"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3c.xlsx  --config data3.cfg data3.csv

# date
## colname typing
../csv2xlsx -d 0 --formatdate "02.01.2006" -o date1.xlsx  date1.csv 
## config file typing, using Go time formatting rules
../csv2xlsx -d 0  --config date1.cfg -o date1b.xlsx  date1b.csv
## use printf + time formats like yyyy-mm-dd
../csv2xlsx -d 0  --config date1b.cfg -o date1bb.xlsx  date1b.csv
## use template
../csv2xlsx -d 0  -o date1a.xlsx  --template template_date.xlsx -r 2 -s Taul1 --headerlines 1 --writeheaderlines 0 date1b.csv

# Benchmark
../csv2xlsx -o out.xlsx data3big.csv  # 4.0 s
../csv2xlsx --font ubuntu --fontsize 12 --verbose 1 -o out.xlsx data3big.csv  # 5.5 s
../csv2xlsx --config data1.cfg -o out.xlsx data3big.csv  # 3.7 s
../csv2xlsx -t template5.xlsx --footer template5footer.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o out.xlsx data3big.csv  # 3.5 s



