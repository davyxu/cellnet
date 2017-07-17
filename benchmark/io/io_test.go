package benchmark

import (
	"testing"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/benchmark"
	_ "github.com/davyxu/cellnet/codec/pb" // 启用pb编码
	"github.com/davyxu/cellnet/proto/pb/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *util.SignalTester

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

	cellnet.RegisterMessage(evd, "gamedef.TestEchoACK", func(ev *cellnet.Event) {

		if qpsm.Acc() > benchmarkSeconds {
			signal.Done(1)
			log.Infof("Average QPS: %d", qpsm.Average())
		}

		ev.Send(&gamedef.TestEchoACK{})

	})

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	evd := socket.NewConnector(queue).Start(benchmarkAddress)

	cellnet.RegisterMessage(evd, "gamedef.TestEchoACK", func(ev *cellnet.Event) {

		ev.Send(&gamedef.TestEchoACK{})

	})

	cellnet.RegisterMessage(evd, "coredef.SessionConnected", func(ev *cellnet.Event) {

		ev.Send(&gamedef.TestEchoACK{})

	})

	queue.StartLoop()

}

func TestIO(t *testing.T) {

	// 屏蔽socket层的调试日志
	golog.SetLevelByString("cellnet", "error")

	signal = util.NewSignalTester(t)

	// 超时时间为测试时间延迟一会
	signal.SetTimeout((benchmarkSeconds + 5) * time.Second)

	server()

	for i := 0; i < clientCount; i++ {
		go client()
	}

	signal.WaitAndExpect("recv time out", 1)

}
