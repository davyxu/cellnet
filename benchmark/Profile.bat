set CURR_DIR=%cd%
cd ../../../../..
set GOPATH=%cd%
go test -v -run=^$ -bench=. -cpuprofile=cpu.pprof -memprofile=mem.pprof github.com/davyxu/cellnet/benchmark/io
: 需要安装Graphviz 
set PATH==%PATH%;"c:\Program Files (x86)\Graphviz2.38\bin"
go tool pprof --pdf io.test.exe cpu.pprof > cpu.pdf
go tool pprof --pdf io.test.exe mem.pprof > mem.pdf
cd %CURR_DIR%