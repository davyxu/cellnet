package gate

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/golang/protobuf/proto"
	"log"
)

var ClientAcceptor cellnet.Peer

// 开启客户端侦听通道
func StartClientAcceptor(pipe cellnet.EventPipe, address string) {

	ClientAcceptor = socket.NewAcceptor(pipe)

	// 所有接收到的消息转发到后台
	ClientAcceptor.InjectData(func(data interface{}) bool {

		if ev, ok := data.(*socket.SessionEvent); ok {

			// Socket各种事件不要往后台发
			switch ev.MsgID {
			case socket.Event_SessionAccepted,
				socket.Event_SessionConnected:
				return true
			}

			// 构建路由封包
			relaypkt := cellnet.BuildPacket(&coredef.UpstreamACK{
				MsgID:    proto.Uint32(ev.MsgID),
				Data:     ev.Data,
				ClientID: proto.Int64(ev.Ses.ID()),
			})

			// TODO 按照封包和逻辑固定分发
			// TODO 非法消息直接掐线
			// TODO 心跳, 超时掐线
			BackendAcceptor.IterateSession(func(ses cellnet.Session) bool {

				if DebugMode {
					log.Printf("client->backend, msgid: %d clientid: %d data: %v", ev.MsgID, ev.Ses.ID(), ev.Data)
				}

				ses.RawSend(relaypkt)

				return true
			})

		}

		return false
	})

	ClientAcceptor.Start(address)

}
