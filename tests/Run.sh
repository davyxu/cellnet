#!/usr/bin/env bash

set -e

# 测试
go test -race -v .

# 出错时，自动删除文件夹
trap 'rm -rf examplebin' EXIT

mkdir -p examplebin

# 保证工程能编译
go build -p 4 -v -o ./examplebin/echo github.com/davyxu/cellnet/examples/echo
go build -p 4 -v -o ./examplebin/echo github.com/davyxu/cellnet/examples/chat/client
go build -p 4 -v -o ./examplebin/echo github.com/davyxu/cellnet/examples/chat/server
go build -p 4 -v -o ./examplebin/echo github.com/davyxu/cellnet/examples/fileserver
go build -p 4 -v -o ./examplebin/echo github.com/davyxu/cellnet/examples/websocket


