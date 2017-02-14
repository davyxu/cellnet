package benchmark

import (
	"testing"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/benchmark"
	"github.com/davyxu/cellnet/example"
	"github.com/davyxu/cellnet/proto/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *test.SignalTester

// 测试地址
const benchmarkAddress = "127.0.0.1:7201"

// 客户端并发数量
const clientCount = 100

// 测试时间(秒)
const benchmarkSeconds = 10

func server() {

	queue := cellnet.NewEventQueue()
	qpsm := benchmark.NewQPSMeter(queue, func(qps int) {

		log.Infof("QPS: %d", qps)

	})

	evd := socket.NewAcceptor(queue).Start(benchmarkAddress)

	socket.RegisterMessage(evd, "gamedef.TestEchoACK", func(ev *cellnet.SessionEvent) {

		if qpsm.Acc() > benchmarkSeconds {
			signal.Done(1)
			log.Infof("Average QPS: %d", qpsm.Average())
		}

		ev.Ses.Send(&gamedef.TestEchoACK{})

	})

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	evd := socket.NewConnector(queue).Start(benchmarkAddress)

	socket.RegisterMessage(evd, "gamedef.TestEchoACK", func(ev *cellnet.SessionEvent) {

		ev.Ses.Send(&gamedef.TestEchoACK{})

	})

	socket.RegisterMessage(evd, "gamedef.SessionConnected", func(ev *cellnet.SessionEvent) {

		ev.Ses.Send(&gamedef.TestEchoACK{})

	})

	queue.StartLoop()

}

func TestIO(t *testing.T) {

	// 屏蔽socket层的调试日志
	golog.SetLevelByString("socket", "error")

	signal = test.NewSignalTester(t)

	// 超时时间为测试时间延迟一会
	signal.SetTimeout((benchmarkSeconds + 5) * time.Second)

	server()

	for i := 0; i < clientCount; i++ {
		go client()
	}

	signal.WaitAndExpect(1, "recv time out")

}
