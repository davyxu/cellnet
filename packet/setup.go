package packet

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

func invokeMsgFunc(ses cellnet.Session, f cellnet.EventFunc, msg interface{}) {
	q := ses.Peer().EventQueue()

	// Peer有队列时，在队列线程调用用户处理函数
	if q != nil {
		q.Post(func() {

			f(msg)
		})

	} else {

		// 在I/O线程调用用户处理函数
		f(msg)
	}
}

func ProcTLVPacket(f cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {
		case socket.RecvEvent: // 接收数据事件
			return onRecvLTVPacket(ev.Ses, f)
		case socket.SendEvent: // 发送数据事件
			return onSendLTVPacket(ev.Ses, f, ev.Msg)
		case socket.ConnectErrorEvent: // 连接错误事件
			invokeMsgFunc(ev.Ses, f, raw)
		case socket.SessionStartEvent: // 会话开始事件（连接上/接受连接）
			invokeMsgFunc(ev.Ses, f, raw)
		case socket.SessionClosedEvent: // 会话关闭事件
			invokeMsgFunc(ev.Ses, f, raw)
		case socket.SessionExitEvent: // 会话退出事件
			invokeMsgFunc(ev.Ses, f, raw)
		case socket.RecvErrorEvent: // 接收错误事件
			log.Errorf("<%s> socket.RecvErrorEvent: %s\n", ev.Ses.Peer().Name(), ev.Error)
			invokeMsgFunc(ev.Ses, f, raw)
		case socket.SendErrorEvent: // 发送错误事件
			log.Errorf("<%s> socket.SendErrorEvent: %s, msg: %#v\n", ev.Ses.Peer().Name(), ev.Error, ev.Msg)
			invokeMsgFunc(ev.Ses, f, raw)
		}

		return nil
	}
}
