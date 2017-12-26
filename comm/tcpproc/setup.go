package tcpproc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/rpc"
)

func ProcLTVInboundPacket(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case *cellnet.ReadEvent: // 接收数据事件

			msg, err := RecvLTVPacket(ev.Ses)
			if err != nil {
				return err
			}

			userFunc(&cellnet.RecvMsgEvent{ev.Ses, msg})
		default:
			userFunc(raw)
		}

		return nil
	}
}

func ProcLTVOutboundPacket(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case *cellnet.SendMsgEvent: // 发送数据事件

			if result := SendLTVPacket(ev.Ses, ev.Msg); result != nil {
				return result
			}
		default:
			userFunc(raw)
		}

		return nil
	}
}

func init() {

	cellnet.RegisterEventProcessor("tcp.ltv", func(userInBound cellnet.EventFunc, userOutbound cellnet.EventFunc) (cellnet.EventFunc, cellnet.EventFunc) {

		return ProcLTVInboundPacket(
				rpc.ProcRPC( // 消息日志
					cellnet.ProcQueue(
						msglog.ProcMsgLog(
							userInBound), // RPC
					),
				),
			),

			msglog.ProcMsgLog(
				rpc.ProcRPC(
					ProcLTVOutboundPacket(userOutbound),
				),
			)

	})
}
