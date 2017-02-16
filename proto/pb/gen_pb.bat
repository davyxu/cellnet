set CURR_DIR=%cd%

: Build generator
cd ..\..\..\..\..\..
set GOPATH=%cd%
go build -o %CURR_DIR%\protoc-gen-msg.exe github.com/davyxu/cellnet/protoc-gen-msg
cd %CURR_DIR%

set outdir=gamedef
set plugindir=..\..\..\..\..\..\bin
mkdir %outdir%
protoc.exe --plugin=protoc-gen-go=%plugindir%\protoc-gen-go.exe --go_out %outdir% --proto_path "." %*
protoc.exe --plugin=protoc-gen-msg=protoc-gen-msg.exe --msg_out=msgid.go:%outdir% --proto_path "." %*