#!/bin/sh


# Basic
../csv2xlsx -d 1 --verbose 1 --font ubuntu --fontsize 20 -o data1.xlsx data1.csv
../csv2xlsx --font ubuntu --fontsize 20 -o data2.xlsx data2.csv
../csv2xlsx -d 1 -o data1.xlsx data1.csv
# Use columnames typing
../csv2xlsx -d 1 -o data4.xlsx data4.csv

# Basic config
../csv2xlsx -d 1 --config data1.cfg --verbose 1 --font ubuntu --fontsize 15 -o data1b.xlsx data1.csv

# Multisheet import
echo "data2 multisheet"
../csv2xlsx -d 1 -c ';' -t template.xlsx -r 4 -s Sh1 -s Sh2 -s NewSheet -o data2.multi.xlsx data2.csv data2.csv data2.csv

echo "data3 using config"
../csv2xlsx -d 0 -c ';' --config data3.cfg  --headerlines 1 --writeheaderlines 1 -s Sh2 -o data3.xlsx  data3.csv

echo "data3 remove header"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3a.xlsx  data3.csv

echo "data3 remove header + footer"
../csv2xlsx -d 1 -c ';' -t template5.xlsx --footer template5footer.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3a.xlsx  data3.csv > x.tmp

echo "data3b not remove header"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 1 -r 5 -s Sh2 -o data3b.xlsx  --verbose 1 data3.csv

echo "data3c remove header, use config"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3c.xlsx  --config data3.cfg data3.csv


