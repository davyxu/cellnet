package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/relay"
	"github.com/davyxu/cellnet/util"
	"reflect"
	"sync"
	"testing"
)

const (
	relayClientToAgent_Address  = "127.0.0.1:16801"
	relayBackendToAgent_Address = "127.0.0.1:16802"

	AgentSessionIDMask = 10000
)

var (
	relay_Signal                  *util.SignalTester
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

			log.Debugln("Relay to agent", relayEvent.Message(), relayEvent.ContextID)
			relay.Relay(relay_BackendToAgentConnector, relayEvent.Message(), relayEvent.ContextID...)

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
	backendToAgentDispatcher := proc.NewMessageDispatcher(relay_BackendToAgentAcceptor, "tcp.ltv")
	backendToAgentDispatcher.RegisterMessage("cellnet.SessionAccepted", func(ev cellnet.Event) {

		backendSession = ev.Session()

		log.Debugln("Backend registed", backendSession.ID())

		wg.Done()

	})

	relay_BackendToAgentAcceptor.Start()

	// 前端侦听
	relay_ClientToAgentAcceptor = peer.NewGenericPeer("tcp.Acceptor", "client->agent", relayClientToAgent_Address, nil)
	ClientToAgentDispatcher := proc.NewMessageDispatcher(relay_ClientToAgentAcceptor, "tcp.ltv")
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

	relay.BindBroadcaster(relay_ClientToAgentAcceptor, relay_BackendToAgentAcceptor, func(frontendPeer cellnet.Peer, ev *relay.RecvMsgEvent) {

		// 广播器
		sesAccessor := frontendPeer.(cellnet.SessionAccessor)

		// 要广播的客户端列表
		for _, maskedSessionID := range ev.ContextID {

			// 去掉掩码
			sesID := maskedSessionID - AgentSessionIDMask
			ses := sesAccessor.GetSession(sesID)
			if ses != nil {

				log.Debugln("Broadcast to client", ev.Message(), sesID)
				ses.Send(ev.Message())
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

	relay_Signal = util.NewSignalTester(t)

	relay_backend()

	relay_agent()

	relay_client()

	relay_Signal.WaitAndExpect("agent not respond", 1)

	relay_ClientToAgentAcceptor.Stop()
	relay_BackendToAgentAcceptor.Stop()
	relay_Client.Stop()

}
