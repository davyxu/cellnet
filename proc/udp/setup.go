package udp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/proc"
)

func ProcLTVInboundPacket(userFunc cellnet.EventProc) cellnet.EventProc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {

		case *cellnet.RecvDataEvent: // 接收数据事件

			msg, err := RecvLTVPacket(ev.Data)
			if err != nil {
				return err
			}

			if _, ok := msg.(*comm.SessionCloseNotify); ok {

				ev.Ses.(udp.UPDSession).RawClose(nil)

			} else {
				userFunc(&cellnet.RecvMsgEvent{ev.Ses, msg})
			}

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
		}

		if userFunc != nil {
			return userFunc(raw)
		}

		return nil
	}
}

func init() {

	proc.RegisterEventProcessor("udp.ltv", func(userInBound cellnet.EventProc, userOutbound cellnet.EventProc) (cellnet.EventProc, cellnet.EventProc) {

		return ProcLTVInboundPacket(
				cellnet.ProcQueue(
					msglog.ProcMsgLog(userInBound),
				),
			),

			msglog.ProcMsgLog(
				ProcLTVOutboundPacket(userOutbound),
			)
	})
}
