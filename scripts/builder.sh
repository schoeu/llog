#!/bin/sh

# linux 64bit
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o ./bin/lla_64bit
upx -9 ./bin/lla_64bit

# linux 32bit
GOOS=linux GOARCH=386 go build -ldflags '-w -s' -o ./bin/lla_32bit
upx -9 ./bin/lla_32bit

# windows 64bit，  windows server暂不考虑
#GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o ./bin/lla_64bit.exe
#upx -9 ./bin/lla_64bit.exe

# windows 32bit，  windows server暂不考虑
#GOOS=windows GOARCH=386 go build -ldflags '-w -s' -o ./bin/lla_32bit.exe
#upx -9 ./bin/lla_32bit.exe

# Mac OS X 64bit，线下测试使用
GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o ./bin/lla_mac
upx -9 ./bin/lla_mac
