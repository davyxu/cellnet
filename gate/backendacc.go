package gate

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"log"
)

var BackendAcceptor cellnet.Peer

// 开启后台服务器的侦听通道
func StartBackendAcceptor(pipe *cellnet.EvPipe, address string) {

	BackendAcceptor = socket.NewAcceptor(pipe)

	// 后台连接开启路由模式
	BackendAcceptor.SetRelayMode(true)

	// 关闭客户端连接
	closeClientMsgID := socket.RegisterSessionMessage(BackendAcceptor, coredef.CloseClientACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.CloseClientACK)

		if msg.ClientID == nil {

			// 关闭所有客户端
			ClientAcceptor.Iterate(func(ses cellnet.Session) bool {

				ses.Close()

				return true
			})

		} else {

			// 关闭指定客户端
			clientSes := ClientAcceptor.Get(msg.GetClientID())

			// 找到连接并关闭
			if clientSes != nil {
				clientSes.Close()
			}

		}

	})

	// 广播
	broardcastMsgID := socket.RegisterSessionMessage(BackendAcceptor, coredef.BroardcastACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.BroardcastACK)

		pkt := &cellnet.Packet{
			MsgID: msg.GetMsgID(),
			Data:  msg.Data,
		}

		if msg.ClientID == nil {

			// 广播给所有客户端
			ClientAcceptor.Iterate(func(ses cellnet.Session) bool {

				ses.(socket.RawSession).RawSend(pkt)

				return true
			})

		} else {

			// 指定客户端发送
			for _, clientid := range msg.ClientID {
				clientSes := ClientAcceptor.Get(clientid)

				if clientSes != nil {
					clientSes.(socket.RawSession).RawSend(pkt)
				}
			}
		}

	})

	// 从后台服务器发来的消息转发到客户端
	BackendAcceptor.Inject(func(data interface{}) bool {

		if ev, ok := data.(*socket.SessionEvent); ok {

			switch ev.MsgID {
			// 收到系统消息时, 透给后方的派发器处理
			case closeClientMsgID,
				broardcastMsgID,
				socket.Event_SessionClosed,
				socket.Event_SessionAccepted:
				return true
			}

			// 找到客户端, 并转发
			clientSes := ClientAcceptor.Get(ev.ClientID)
			if clientSes != nil {

				if DebugMode {
					log.Printf("backend->client, msgid: %d clientid %d", ev.MsgID, ev.ClientID)
				}
				clientSes.(socket.RawSession).RawSend(ev.Packet)

			} else if DebugMode {

				log.Printf("backend->client, client not found, msgid: %d clientid %d", ev.MsgID, ev.ClientID)
			}

		}

		return false
	})

	BackendAcceptor.Start(address)
}
