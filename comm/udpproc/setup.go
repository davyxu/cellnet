package udpproc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
)

func ProcLTVInboundPacket(userFunc cellnet.EventProc) cellnet.EventProc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {

		case *cellnet.RecvDataEvent: // 接收数据事件

			msg, err := RecvLTVPacket(ev.Data)
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
		}

		if userFunc != nil {
			return userFunc(raw)
		}

		return nil
	}
}

func init() {

	cellnet.RegisterEventProcessor("udp.ltv", func(userInBound cellnet.EventProc, userOutbound cellnet.EventProc) (cellnet.EventProc, cellnet.EventProc) {

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
