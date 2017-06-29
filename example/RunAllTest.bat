set CURRDIR=%cd%
cd ../../../../..
set GOPATH=%cd%

go test github.com/davyxu/cellnet/example/classicrecv ^
github.com/davyxu/cellnet/example/classicrecv ^
github.com/davyxu/cellnet/example/close ^
github.com/davyxu/cellnet/example/echo_pb ^
github.com/davyxu/cellnet/example/echo_sproto ^
github.com/davyxu/cellnet/example/rpc ^
github.com/davyxu/cellnet/example/timer

cd %CURRDIR%