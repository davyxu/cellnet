#!/usr/bin/env bash

# 设置GOPATH
CURR=`pwd`
cd ../../../../../../..
export GOPATH=`pwd`
cd ${CURR}

# 插件及protoc存放路径
BIN_PATH=${GOPATH}/bin

# 错误时退出
set -e

go install -v github.com/gogo/protobuf/protoc-gen-gogofaster

go install -v github.com/davyxu/cellnet/protoc-gen-msg

# windows下时，添加后缀名
if [ `go env GOHOSTOS` == "windows" ];then
	EXESUFFIX=.exe
fi

# 生成协议
${BIN_PATH}/protoc --plugin=protoc-gen-gogofaster=${BIN_PATH}/protoc-gen-gogofaster${EXESUFFIX} --gogofaster_out=. --proto_path="." pb.proto

# 生成cellnet 消息注册文件
${BIN_PATH}/protoc --plugin=protoc-gen-msg=${BIN_PATH}/protoc-gen-msg${EXESUFFIX} --msg_out=msgid.go:. --proto_path="." pb.proto
