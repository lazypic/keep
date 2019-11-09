#!/bin/sh
APP="keep"
GOOS=linux GOARCH=amd64 go build -o ./bin/linux/${APP} main.go
GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/${APP} main.go
GOOS=windows GOARCH=amd64 go build -o ./bin/windows/${APP} main.go

# Github Release에 업로드 하기위해 압축
cd ./bin/linux/ && tar -zcvf ../${APP}_linux_x86-64.tgz . && cd -
cd ./bin/darwin/ && tar -zcvf ../${APP}_darwin_x86-64.tgz . && cd -
cd ./bin/windows/ && tar -zcvf ../${APP}_windows_x86-64.tgz . && cd -

# 삭제
rm -rf ./bin/linux
rm -rf ./bin/darwin
rm -rf ./bin/windows
