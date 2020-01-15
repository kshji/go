#!/bin/sh

ver=2020-01-09

env GOOS=linux go build -o "csv2xlsx.linux.${ver}" ../
env GOOS=darwin go build -o "csv2xlsx.darwin.${ver}" ../
env GOOS=windows go build -o "csv2xlsx.win.${ver}.exe" ../

cp "csv2xlsx.linux.${ver}" csv2xlsx.linux
cp "csv2xlsx.darwin.${ver}" csv2xlsx.osx
cp "csv2xlsx.win.${ver}.exe" csv2xlsx.exe
