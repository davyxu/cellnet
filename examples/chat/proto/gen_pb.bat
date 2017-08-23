set CURR_DIR=%cd%

: Build generator
cd ..\..\..\..\..\..\..
set GOPATH=%cd%
go build -o %CURR_DIR%\protoc-gen-msg.exe github.com/davyxu/cellnet/protoc-gen-msg
go build -o %CURR_DIR%\protoc-gen-go.exe github.com/golang/protobuf/protoc-gen-go
cd %CURR_DIR%

set OUTDIR=%PKGNAME%
mkdir %OUTDIR%
protoc.exe --plugin=protoc-gen-go=protoc-gen-go.exe --go_out %OUTDIR% --proto_path "." %*
protoc.exe --plugin=protoc-gen-msg=protoc-gen-msg.exe --msg_out=msgid.go:%OUTDIR% --proto_path "." %*