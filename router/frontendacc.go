package router

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
)

// agent服务器处理客户端连接时使用以下代码

var FrontendAcceptor cellnet.Peer

func getMsgName(msgid uint32) string {

	if meta := cellnet.MessageMetaByID(msgid); meta != nil {
		return meta.Name
	}

	return ""
}

// 开启前端侦听通道
func StartFrontendAcceptor(pipe cellnet.EventPipe, address string, peerName string) {

	FrontendAcceptor = socket.NewAcceptor(pipe)
	FrontendAcceptor.SetName(peerName)

	// 默认开启并发
	FrontendAcceptor.EnableConcurrenceMode(true)

	// 所有接收到的消息转发到后台
	FrontendAcceptor.InjectData(func(data interface{}) bool {

		if ev, ok := data.(*socket.SessionEvent); ok {

			// Socket各种事件不要往后台发
			switch ev.MsgID {
			case socket.Event_SessionAccepted,
				socket.Event_SessionConnected:
				return true
			}

			// TODO 非法消息直接掐线
			// TODO 心跳, 超时掐线

			// 广播到后台所有服务器
			if relayMethod == RelayMethod_BroardcastToAllBackend || ev.MsgID == socket.Event_SessionClosed {

				BackendAcceptor.IterateSession(func(ses cellnet.Session) bool {

					sendMessageToBackend(ses, ev)

					return true
				})

				// 按照白名单准确投递
			} else if relayMethod == RelayMethod_WhiteList {
				ses := getRelaySession(ev.MsgID)

				if ses != nil {

					sendMessageToBackend(ses, ev)

				} else {

					if DebugMode {
						log.Errorf("client->backend, msg: %s(%d) clientid: %d  relay target not found", getMsgName(ev.MsgID), ev.MsgID, ev.Ses.ID())
					}

				}

			}

		}

		return false
	})

	FrontendAcceptor.Start(address)

}

func sendMessageToBackend(ses cellnet.Session, ev *socket.SessionEvent) {

	// 构建路由封包
	relaypkt, _ := cellnet.BuildPacket(&coredef.UpstreamACK{
		MsgID:    ev.MsgID,
		Data:     ev.Data,
		ClientID: ev.Ses.ID(),
	})

	if DebugMode {
		log.Debugf("client->backend, msg: %s(%d) clientid: %d", getMsgName(ev.MsgID), ev.MsgID, ev.Ses.ID())
	}

	ses.RawSend(relaypkt)
}
