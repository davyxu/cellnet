package benchmark

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/benchmark"
	"github.com/davyxu/cellnet/packet"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/tests/proto"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
	"testing"
	"time"
)

var log = golog.New("test")

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

	var config cellnet.PeerConfig
	config.TypeName = "tcp.Acceptor"
	config.Queue = queue
	config.Address = benchmarkAddress
	config.Name = "server"
	config.InboundEvent = packet.ProcTLVPacket(func(ses cellnet.Session, raw interface{}) {

		switch raw.(type) {
		case packet.RecvMsgEvent:

			if qpsm.Acc() > benchmarkSeconds {
				signal.Done(1)
				log.Infof("Average QPS: %d", qpsm.Average())
			}

			ses.Send(&proto.TestEchoACK{})
		}

	})

	cellnet.NewPeer(config).Start()

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	var config cellnet.PeerConfig
	config.TypeName = "tcp.Connector"
	config.Queue = queue
	config.Address = benchmarkAddress
	config.Name = "client"
	config.InboundEvent = packet.ProcTLVPacket(func(ses cellnet.Session, raw interface{}) {

		switch raw.(type) {
		case socket.ConnectedEvent:
			ses.Send(&proto.TestEchoACK{})
		case packet.RecvMsgEvent:

			ses.Send(&proto.TestEchoACK{})
		}

	})

	cellnet.NewPeer(config).Start()

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
