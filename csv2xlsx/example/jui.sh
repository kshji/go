#!/bin/sh

#../csv2xlsx -d 1 -c ';' -t template.xlsx -r 4 -s Sh1 -s Sh2 -s NewSheet -o out.xlsx data2.csv data2.csv data2.csv
testi()
{
echo "data2"
../csv2xlsx -d 0 -c ';' -t template2.xlsx --headerlines 0 -r 4 -s Sh2 -o data2.xlsx data2.csv 
echo "data3"
../csv2xlsx -d 0 -c ';' -t template3.xlsx --headerlines 1 -r 5 -s Sh2 -o data3.xlsx data3.csv 
echo "data3b"
../csv2xlsx -d 0 -c ';' -t template3.xlsx --writeheaderlines 0 --startrow 2 -r 5 -s Sh2 -o data3b.xlsx data3.csv 
echo "data3c"
../csv2xlsx -d 0 -c ';' -t template3.xlsx --writeheaderlines 0 --startrow 2 -r 5 -s Sh2 -o data3c.xlsx data3.csv 

}

#MAIN


echo "data3 remove header"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data3.xlsx  data3.csv
echo "data3b not remove header"
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 1 -r 5 -s Sh2 -o data3b.xlsx  data3.csv
echo "data4 total "
../csv2xlsx -d 0 -c ';' -t template5.xlsx --headerlines 1 --writeheaderlines 0 -r 5 -s Sh2 -o data4.xlsx  data4.csv

echo "big.xlsx"
../csv2xlsx -d 0 -c ';' -t varasto.xlsx -o big.xlsx -r 3 --headerlines 1 --writeheaderlines 0 -s Varasto --verbose 1 big.csv
