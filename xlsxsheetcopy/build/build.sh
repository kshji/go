#!/bin/sh

ver=2019-12-26
prg=xlsxsheetcopy

env GOOS=linux go build -o "$prg.linux.${ver}" ../
env GOOS=darwin go build -o "$prg.darwin.${ver}" ../
env GOOS=windows go build -o "$prg.win.${ver}.exe" ../

cp "$prg.linux.${ver}" $prg.linux
cp "$prg.darwin.${ver}" $prg.osx
cp "$prg.win.${ver}.exe" $prg.exe
