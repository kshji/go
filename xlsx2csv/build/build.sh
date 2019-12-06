#!/bin/sh

ver=2019-12-06

env GOOS=linux go build -o "xlsx2csv.linux.${ver}" ../
env GOOS=darwin go build -o "xlsx2csv.mac.${ver}" ../
env GOOS=windows go build -o "xlsx2csv.win.${ver}.exe" ../
