package tcpproc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm/rpc"
	"github.com/davyxu/cellnet/msglog"
)

func ProcLTVInboundPacket(userFunc cellnet.EventProc) cellnet.EventProc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case *cellnet.ReadStreamEvent: // 接收数据事件

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

func ProcLTVOutboundPacket(userFunc cellnet.EventProc) cellnet.EventProc {

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

	cellnet.RegisterEventProcessor("tcp.ltv", func(userInBound cellnet.EventProc, userOutbound cellnet.EventProc) (cellnet.EventProc, cellnet.EventProc) {

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
