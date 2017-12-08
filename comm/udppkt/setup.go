package udppkt

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
)

func ProcTVPacket(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {

		case cellnet.RecvDataEvent: // 接收数据事件

			if result := onRecvTVPacket(ev.Ses, ev.Data, userFunc); result != nil {
				return result
			}

		case cellnet.SendMsgEvent: // 发送数据事件

			if result := onSendTVPacket(ev.Ses, ev.Msg); result != nil {
				return result
			}
		}

		return userFunc(raw)
	}
}
func ProcQueue(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case cellnet.RecvMsgEvent:

			cellnet.QueuedCall(ev.Ses, func() {
				userFunc(raw)
			})

		default:
			return userFunc(raw)
		}

		return nil
	}
}

func ProcSysMsg(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		if userFunc == nil {
			return nil
		}

		switch ev := raw.(type) {

		case cellnet.SessionConnectErrorEvent:
			userFunc(cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionConnectError{}})
		case cellnet.SessionClosedEvent:
			userFunc(cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionClosed{}})
		case cellnet.SessionAcceptedEvent:
			userFunc(cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionAccepted{}})
		case cellnet.SessionConnectedEvent:
			userFunc(cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionConnected{}})
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

	return ProcTVPacket(
		msglog.ProcMsgLog(
			ProcSysMsg(final),
		),
	)
}

func init() {

	cellnet.RegisterPeerCreator("tv.udp.Connector", func(config cellnet.PeerConfig) cellnet.Peer {

		config.PeerType = "udp.Connector"
		p := cellnet.NewPeer(config)

		p.SetEventFunc(initEvent(&config))

		return p
	})

	cellnet.RegisterPeerCreator("tv.udp.Acceptor", func(config cellnet.PeerConfig) cellnet.Peer {
		config.PeerType = "udp.Acceptor"
		p := cellnet.NewPeer(config)

		p.SetEventFunc(initEvent(&config))

		return p
	})
}
