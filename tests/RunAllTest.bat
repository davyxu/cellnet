set CURRDIR=%cd%
cd ../../../../..
set GOPATH=%cd%
cd %CURRDIR%

go test -v .


