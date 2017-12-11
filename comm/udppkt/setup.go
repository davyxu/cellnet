package udppkt

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/msglog"
)

func ProcTVPacket(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {

		case cellnet.RecvDataEvent: // 接收数据事件

			if result := onRecvLTVPacket(ev.Ses, ev.Data, userFunc); result != nil {
				return result
			}

		case cellnet.SendMsgEvent: // 发送数据事件

			if result := onSendLTVPacket(ev.Ses, ev.Msg); result != nil {
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
			comm.ProcSysMsg(final),
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
