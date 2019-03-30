package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/relay"
	"reflect"
	"sync"
	"testing"
	"time"
)

const (
	relayClientToAgent_Address  = "127.0.0.1:16801"
	relayBackendToAgent_Address = "127.0.0.1:16802"

	AgentSessionIDMask = 10000
)

var (
	relay_Signal                  *SignalTester
	relay_ClientToAgentAcceptor   cellnet.Peer
	relay_BackendToAgentConnector cellnet.Peer
	relay_BackendToAgentAcceptor  cellnet.Peer
	relay_Client                  cellnet.Peer
)

func relay_backend() {
	queue := cellnet.NewEventQueue()

	relay_BackendToAgentConnector = peer.NewGenericPeer("tcp.Connector", "backend", relayBackendToAgent_Address, queue)

	proc.BindProcessorHandler(relay_BackendToAgentConnector, "tcp.ltv", func(ev cellnet.Event) {

		if relayEvent, ok := ev.(*relay.RecvMsgEvent); ok {

			log.Debugln("Relay to agent", relayEvent.Message(), relayEvent.PassThroughAsInt64())
			relay.Relay(relay_BackendToAgentConnector, relayEvent.Message(), relayEvent.PassThroughAsInt64())

		}

	})

	relay_BackendToAgentConnector.Start()

	queue.StartLoop()
}

func relay_agent() {

	var backendSession cellnet.Session

	var wg sync.WaitGroup
	wg.Add(1)
	// 后端侦听
	relay_BackendToAgentAcceptor = peer.NewGenericPeer("tcp.Acceptor", "backend->agent", relayBackendToAgent_Address, nil)
	backendToAgentDispatcher := proc.NewMessageDispatcherBindPeer(relay_BackendToAgentAcceptor, "tcp.ltv")
	backendToAgentDispatcher.RegisterMessage("cellnet.SessionAccepted", func(ev cellnet.Event) {

		backendSession = ev.Session()

		log.Debugln("Backend registered", backendSession.ID())

		wg.Done()

	})

	relay_BackendToAgentAcceptor.Start()

	// 前端侦听
	relay_ClientToAgentAcceptor = peer.NewGenericPeer("tcp.Acceptor", "client->agent", relayClientToAgent_Address, nil)
	ClientToAgentDispatcher := proc.NewMessageDispatcherBindPeer(relay_ClientToAgentAcceptor, "tcp.ltv")
	ClientToAgentDispatcher.RegisterMessage("tests.TestEchoACK", func(ev cellnet.Event) {

		// 等待后台会话连接后，再转发消息给后台，本Test专用
		wg.Wait()

		// 只有在后端服务器连接时
		if backendSession != nil {

			// 添加掩码的sesid
			maskedSessionID := ev.Session().ID() + AgentSessionIDMask

			log.Debugln("Relay to backend", ev.Message(), ev.Session().ID())
			// 路由到后台
			relay.Relay(backendSession, ev.Message(), maskedSessionID)
		} else {
			panic("backendSession is not ready")
		}

	})

	relay_ClientToAgentAcceptor.Start()

	relay.SetBroadcaster(func(event *relay.RecvMsgEvent) {

		// 仅限于从后端来的Relay消息, 本Test中，因为3个进程逻辑混在一起，必须这样区分来源
		if event.Ses.Peer() == relay_BackendToAgentAcceptor {

			// 广播器
			sesAccessor := relay_ClientToAgentAcceptor.(cellnet.SessionAccessor)

			// 去掉掩码
			sesID := event.PassThroughAsInt64() - AgentSessionIDMask
			ses := sesAccessor.GetSession(sesID)
			if ses != nil {

				log.Debugln("Broadcast to client", event.Message(), sesID)
				ses.Send(event.Message())
			}
		}

	})

}

func relay_client() {

	queue := cellnet.NewEventQueue()

	relay_Client = peer.NewGenericPeer("tcp.Connector", "client", relayClientToAgent_Address, queue)

	dataMsg := TestEchoACK{
		Msg:   "hello",
		Value: 1234,
	}

	proc.BindProcessorHandler(relay_Client, "tcp.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:
			log.Debugln("send data")
			ev.Session().Send(&dataMsg)

		case *TestEchoACK:
			if reflect.DeepEqual(dataMsg, *msg) {
				relay_Signal.Done(1)
				log.Debugln("data done")
			}
		}

	})

	relay_Client.Start()

	queue.StartLoop()

}

func TestRelay(t *testing.T) {

	relay_Signal = NewSignalTester(t)
	relay_Signal.SetTimeout(time.Second * 5)

	relay_agent()

	relay_backend()

	relay_client()

	relay_Signal.WaitAndExpect("agent not respond", 1)

	relay_ClientToAgentAcceptor.Stop()
	relay_BackendToAgentAcceptor.Stop()
	relay_Client.Stop()

}
