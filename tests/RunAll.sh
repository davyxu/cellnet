#!/usr/bin/env bash
export CURRDIR=`pwd`
cd ../../../../..
export GOPATH=`pwd`
cd ${CURRDIR}

go test -v .
if [ $? -ne 0 ] ; then read -rsp $'Errors occurred...\n' ; fi