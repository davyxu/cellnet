set CURR_DIR=%cd%

: Build generator
cd ..\..\..\..\..\..
set GOPATH=%cd%
go build -o %CURR_DIR%\sprotogen.exe github.com/davyxu/gosproto/sprotogen
cd %CURR_DIR%

: Generate go source file by sproto
sprotogen --go_out=.\gamedef\gamedef.go --package=gamedef --cellnet_reg=true %*