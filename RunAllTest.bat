set CURRDIR=%cd%
cd ../../../..
set GOPATH=%cd%

go test github.com/davyxu/cellnet/example/classicrecv ^
github.com/davyxu/cellnet/example/sendclose ^
github.com/davyxu/cellnet/example/echo_pb ^
github.com/davyxu/cellnet/example/echo_sproto ^
github.com/davyxu/cellnet/example/echo_websocket ^
github.com/davyxu/cellnet/example/gracefulexit ^
github.com/davyxu/cellnet/example/rpc ^
github.com/davyxu/cellnet/example/timer


cd %CURRDIR%