package msglog

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
)

func ProcMsgLog(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case *cellnet.RecvMsgEvent:

			if IsBlockedMessageByID(cellnet.MessageID(ev.Msg)) {
				break
			}

			switch ev.Msg.(type) {
			case *comm.SessionAccepted:
				log.Debugf("#accepted(%s)@%d", ev.Ses.Peer().Name(), ev.Ses.ID())
			case *comm.SessionClosed:
				log.Debugf("#closed(%s)@%d", ev.Ses.Peer().Name(), ev.Ses.ID())
			case *comm.SessionConnected:
				log.Debugf("#connected(%s)@%d", ev.Ses.Peer().Name(), ev.Ses.ID())
			case *comm.SessionConnectError:
				log.Debugf("#connectfailed(%s)@%d address: %s", ev.Ses.Peer().Name(), ev.Ses.ID(), ev.Ses.Peer().Address())
			default:
				log.Debugf("#recv(%s)@%d %s(%d) | %s",
					ev.Ses.Peer().Name(),
					ev.Ses.ID(),
					cellnet.MessageName(ev.Msg),
					cellnet.MessageID(ev.Msg),
					cellnet.MessageToString(ev.Msg))
			}

		case *cellnet.SendMsgEvent:

			if IsBlockedMessageByID(cellnet.MessageID(ev.Msg)) {
				break
			}

			log.Debugf("#send(%s)@%d %s(%d) | %s",
				ev.Ses.Peer().Name(),
				ev.Ses.ID(),
				cellnet.MessageName(ev.Msg),
				cellnet.MessageID(ev.Msg),
				cellnet.MessageToString(ev.Msg))

		case *cellnet.RecvErrorEvent: // 接收错误事件
			log.Debugf("#recverror(%s)@%d address: %s, %s", ev.Ses.Peer().Name(), ev.Ses.ID(), ev.Ses.Peer().Address(), ev.Error)
		case *cellnet.SendMsgErrorEvent: // 发送错误事件
			log.Debugf("#senderror(%s)@%d address: %s, %s", ev.Ses.Peer().Name(), ev.Ses.ID(), ev.Ses.Peer().Address(), ev.Error)
		}

		if userFunc != nil {
			return userFunc(raw)
		}

		return nil
	}
}
