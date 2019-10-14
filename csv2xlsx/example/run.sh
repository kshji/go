#!/bin/sh

# Multisheet import
echo "data2 multisheet"
../csv2xlsx -d 1 -c ';' -t template.xlsx -r 4 -s Sh1 -s Sh2 -s NewSheet -o out.xlsx data2.csv data2.csv data2.csv

echo "data3 remove header"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3.xlsx  data3.csv

echo "data3b not remove header"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 1 -r 5 -s Sh2 -o data3b.xlsx  --verbose 1 data3.csv


