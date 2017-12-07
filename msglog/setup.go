package msglog

import (
	"github.com/davyxu/cellnet"
)

func ProcMsgLog(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case cellnet.SessionConnectErrorEvent: // 连接错误事件
			log.Debugf("#connectfailed(%s)@%d address: %s", ev.Ses.Peer().Name(), ev.Ses.ID(), ev.Ses.Peer().Address())
		case cellnet.SessionConnectedEvent: // 会话开始事件（连接上/接受连接）
			log.Debugf("#connected(%s)@%d", ev.Ses.Peer().Name(), ev.Ses.ID())
		case cellnet.SessionAcceptedEvent:

			log.Debugf("#accepted(%s)@%d", ev.Ses.Peer().Name(), ev.Ses.ID())

		case cellnet.SessionClosedEvent: // 会话关闭事件
			log.Debugf("#closed(%s)@%d", ev.Ses.Peer().Name(), ev.Ses.ID())

		case cellnet.RecvMsgEvent:

			if IsBlockedMessageByID(MessageID(ev.Msg)) {
				break
			}

			log.Debugf("#recv(%s)@%d %s(%d) | %s",
				ev.Ses.Peer().Name(),
				ev.Ses.ID(),
				MessageName(ev.Msg),
				MessageID(ev.Msg),
				MessageToString(ev.Msg))

		case cellnet.SendMsgEvent:

			if IsBlockedMessageByID(MessageID(ev.Msg)) {
				break
			}

			log.Debugf("#send(%s)@%d %s(%d) | %s",
				ev.Ses.Peer().Name(),
				ev.Ses.ID(),
				MessageName(ev.Msg),
				MessageID(ev.Msg),
				MessageToString(ev.Msg))

		case cellnet.RecvErrorEvent: // 接收错误事件
			log.Debugf("#recverror(%s)@%d address: %s, %s", ev.Ses.Peer().Name(), ev.Ses.ID(), ev.Ses.Peer().Address(), ev.Error)
		case cellnet.SendMsgErrorEvent: // 发送错误事件
			log.Debugf("#senderror(%s)@%d address: %s, %s", ev.Ses.Peer().Name(), ev.Ses.ID(), ev.Ses.Peer().Address(), ev.Error)
		}

		return userFunc(raw)
	}
}
