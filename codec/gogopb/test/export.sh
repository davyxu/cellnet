#!/usr/bin/env bash
CURR=`pwd`
cd ../../../../../../..
export GOPATH=`pwd`
cd ${CURR}

go build -v -o=protoc-gen-gogofaster github.com/gogo/protobuf/protoc-gen-gogofaster

go build -v -o=protoc-gen-msg github.com/davyxu/cellnet/protoc-gen-msg

# 生成协议
./protoc --plugin=protoc-gen-gogofaster=protoc-gen-gogofaster --gogofaster_out=. --proto_path="." pb.proto
if [ $? -ne 0 ] ; then read -rsp $'Errors occurred...\n' ; fi

# 生成cellnet 消息注册文件
./protoc --plugin=protoc-gen-msg=protoc-gen-msg --msg_out=msgid.go:. --proto_path="." pb.proto
if [ $? -ne 0 ] ; then read -rsp $'Errors occurred...\n' ; fi
