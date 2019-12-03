#!/bin/sh

ver=2019-12-02

env GOOS=linux go build -o "csv2xlsx.linux.${ver}" ../
env GOOS=darwin go build -o "csv2xlsx.darwin.${ver}" ../
env GOOS=windows go build -o "csv2xlsx.win.${ver}.exe" ../
