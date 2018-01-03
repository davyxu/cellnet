package kcpproc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/msglog"
)

const kcpTag = "kcp"

func mustKCPContext(ses cellnet.Session) *kcpContext {
	if kcpRaw, ok := ses.GetTag(kcpTag); ok {
		return kcpRaw.(*kcpContext)
	} else {
		panic("invalid kcp context")
	}
}

func ProcKCPInboundPacket(userFunc cellnet.EventProc) cellnet.EventProc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case *cellnet.RecvMsgEvent:

			switch ev.Msg.(type) {
			case *comm.SessionAccepted,
				*comm.SessionConnected:
				ev.Ses.SetTag(kcpTag, newContext(ev.Ses, userFunc))
			case *comm.SessionClosed:
				mustKCPContext(ev.Ses).Close()
			}

			userFunc(raw)

		case *cellnet.RecvDataEvent: // 接收数据事件

			mustKCPContext(ev.Ses).input(ev.Data)

		default:
			userFunc(raw)
		}

		return nil
	}
}

func ProcKCPOutboundPacket(userFunc cellnet.EventProc) cellnet.EventProc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case *cellnet.SendMsgEvent: // 发送数据事件

			if result := mustKCPContext(ev.Ses).sendMessage(ev.Msg); result != nil {
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

	cellnet.RegisterEventProcessor("udp.kcp", func(userInBound cellnet.EventProc, userOutbound cellnet.EventProc) (cellnet.EventProc, cellnet.EventProc) {

		return ProcKCPInboundPacket(
				cellnet.ProcQueue(
					msglog.ProcMsgLog(userInBound),
				),
			),

			msglog.ProcMsgLog(
				ProcKCPOutboundPacket(userOutbound),
			)
	})
}
