package packet

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

func ProcQueued(f cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {
		switch ev := raw.(type) {
		case socket.RecvEvent: // 接收数据事件
			cellnet.CallEventFuncAutoQueue(ev.Ses, f, raw)
		case socket.SendEvent: // 发送数据事件
			cellnet.CallEventFuncAutoQueue(ev.Ses, f, raw)
		case socket.ConnectErrorEvent: // 连接错误事件
			cellnet.CallEventFuncAutoQueue(ev.Ses, f, raw)
		case socket.SessionStartEvent: // 会话开始事件（连接上/接受连接）
			cellnet.CallEventFuncAutoQueue(ev.Ses, f, raw)
		case socket.SessionClosedEvent: // 会话关闭事件
			cellnet.CallEventFuncAutoQueue(ev.Ses, f, raw)
		case socket.SessionExitEvent: // 会话退出事件
			cellnet.CallEventFuncAutoQueue(ev.Ses, f, raw)
		case socket.RecvErrorEvent: // 接收错误事件
			cellnet.CallEventFuncAutoQueue(ev.Ses, f, raw)
		case socket.SendErrorEvent: // 发送错误事件
			cellnet.CallEventFuncAutoQueue(ev.Ses, f, raw)
		}

		return nil
	}
}

func ProcTLVPacket(f cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {

		case socket.RecvEvent: // 接收数据事件
			return onRecvLTVPacket(ev.Ses, f)
		case socket.SendEvent: // 发送数据事件
			return onSendLTVPacket(ev.Ses, f, ev.Msg)
		}

		return f(raw)
	}
}
