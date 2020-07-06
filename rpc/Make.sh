#!/usr/bin/env bash
CURRDIR=`pwd`
cd ../../../../..
export GOPATH=`pwd`
cd ${CURRDIR}

go build -v -o=${GOPATH}/bin/protoplus github.com/davyxu/protoplus/cmd/protoplus


${GOPATH}/bin/protoplus -go_out=msg_gen.go -package=rpc -genreg msg.proto