package msglog

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/packet"
	"github.com/davyxu/cellnet/socket"
	"reflect"
	"strings"
)

func MessageName(msg interface{}) string {

	meta := cellnet.MessageMetaByType(reflect.TypeOf(msg).Elem())
	if meta == nil {
		return ""
	}

	return meta.Name
}

func MessageToString(msg interface{}) string {

	if msg == nil {
		return ""
	}

	if stringer, ok := msg.(interface {
		String() string
	}); ok {
		return stringer.String()
	}

	return ""
}

func ProcMsgLog(f cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case socket.ConnectErrorEvent: // 连接错误事件
			log.Debugf("#connectfailed(%s)@%d address: %s", ev.Ses.Peer().Name(), ev.Ses.ID(), ev.Ses.Peer().Address())
		case socket.SessionStartEvent: // 会话开始事件（连接上/接受连接）

			if strings.Contains(ev.Ses.Peer().TypeName(), "Acceptor") {
				log.Debugf("#accepted(%s)@%d", ev.Ses.Peer().Name(), ev.Ses.ID())
			} else if strings.Contains(ev.Ses.Peer().TypeName(), "Connector") {
				log.Debugf("#connected(%s)@%d", ev.Ses.Peer().Name(), ev.Ses.ID())
			}

		case socket.SessionClosedEvent: // 会话关闭事件
			log.Debugf("#closed(%s)@%d", ev.Ses.Peer().Name(), ev.Ses.ID())
		case socket.SessionExitEvent: // 会话退出事件

		case packet.RecvMsgEvent:

			log.Debugf("#recv(%s)@%d %s(%d) size: %d | %s",
				ev.Ses.Peer().Name(),
				ev.Ses.ID(),
				MessageName(ev.Msg),
				ev.MsgID,
				len(ev.MsgData), MessageToString(ev.Msg))

		case packet.SendMsgEvent:

			data, msgid, _ := cellnet.EncodeMessage(ev.Msg)

			log.Debugf("#send(%s)@%d %s(%d) size: %d | %s",
				ev.Ses.Peer().Name(),
				ev.Ses.ID(),
				MessageName(ev.Msg),
				msgid,
				len(data), MessageToString(ev.Msg))

		case socket.RecvErrorEvent: // 接收错误事件
			log.Debugf("#recverror(%s)@%d address: %s, %s", ev.Ses.Peer().Name(), ev.Ses.ID(), ev.Ses.Peer().Address(), ev.Error)
		case socket.SendErrorEvent: // 发送错误事件
			log.Debugf("#senderror(%s)@%d address: %s, %s", ev.Ses.Peer().Name(), ev.Ses.ID(), ev.Ses.Peer().Address(), ev.Error)
		}

		return f(raw)
	}
}
