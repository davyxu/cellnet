package tcppkt

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/rpc"
)

func ProcQueue(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case *cellnet.RecvMsgEvent:

			cellnet.QueuedCall(ev.Ses, func() {
				userFunc(raw)
			})

		case *rpc.RecvMsgEvent:

			q := ev.Queue()
			// Peer有队列时，在队列线程调用用户处理函数
			if q != nil {
				q.Post(func() {
					userFunc(raw)
				})

			} else {

				// 在I/O线程调用用户处理函数
				return userFunc(raw)
			}
		default:
			return userFunc(raw)
		}

		return nil
	}
}

func ProcTLVPacket(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {

		case *cellnet.ReadEvent: // 接收数据事件

			if result := onRecvLTVPacket(ev.Ses, userFunc); result != nil {
				return result
			}

		case *cellnet.SendMsgEvent: // 发送数据事件

			if result := onSendLTVPacket(ev.Ses, ev.Msg); result != nil {
				return result
			}
		}

		return userFunc(raw)
	}
}

func initEvent(config *cellnet.PeerConfig) cellnet.EventFunc {
	var final cellnet.EventFunc

	// 有队列，添加队列处理
	if config.Queue != nil {
		final = ProcQueue(config.Event)
	} else {
		// 否则直接处理
		final = config.Event
	}

	return ProcTLVPacket(
		rpc.ProcRPC( // 消息日志
			comm.ProcSysMsg( // 系统事件转消息
				msglog.ProcMsgLog(final), // RPC
			),
		),
	)
}

func init() {

	cellnet.RegisterPeerCreator("ltv.tcp.Connector", func(config cellnet.PeerConfig) cellnet.Peer {

		config.PeerType = "tcp.Connector"
		p := cellnet.NewPeer(config)

		p.SetEventFunc(initEvent(&config))

		return p
	})

	cellnet.RegisterPeerCreator("ltv.tcp.Acceptor", func(config cellnet.PeerConfig) cellnet.Peer {
		config.PeerType = "tcp.Acceptor"
		p := cellnet.NewPeer(config)

		p.SetEventFunc(initEvent(&config))

		return p
	})
}
