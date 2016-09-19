package benchmark

import (
	"testing"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/benchmark"
	"github.com/davyxu/cellnet/proto/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/test"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *test.SignalTester

// 测试地址
const benchmarkAddress = "127.0.0.1:7201"

// 客户端并发数量
const clientCount = 1000

// 测试时间(秒)
const benchmarkSeconds = 10

func server() {

	pipe := cellnet.NewEventPipe()

	qpsm := benchmark.NewQPSMeter(pipe, func(qps int) {

		log.Infof("QPS: %d", qps)

	})

	evq := socket.NewAcceptor(pipe).Start(benchmarkAddress)

	socket.RegisterSessionMessage(evq, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {

		if qpsm.Acc() > benchmarkSeconds {
			signal.Done(1)
			log.Infof("Average QPS: %d", qpsm.Average())
		}

		ses.Send(&gamedef.TestEchoACK{})

	})

	pipe.Start()

}

func client() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewConnector(pipe).Start(benchmarkAddress)

	socket.RegisterSessionMessage(evq, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {

		ses.Send(&gamedef.TestEchoACK{})

	})

	socket.RegisterSessionMessage(evq, "gamedef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		ses.Send(&gamedef.TestEchoACK{})

	})

	pipe.Start()

}

func TestBenchmark(t *testing.T) {

	// 屏蔽socket层的调试日志
	golog.SetLevelByString("socket", "info")

	signal = test.NewSignalTester(t)

	// 超时时间为测试时间延迟一会
	signal.SetTimeout((benchmarkSeconds + 5) * time.Second)

	server()

	for i := 0; i < clientCount; i++ {
		go client()
	}

	signal.WaitAndExpect(1, "recv time out")

}
