package gate

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
	"log"
)

var ClientAcceptor cellnet.Peer

// 开启客户端侦听通道
func StartClientAcceptor(pipe *cellnet.EvPipe, address string) {

	ClientAcceptor = socket.NewAcceptor(pipe)

	// 所有接收到的消息转发到后台
	ClientAcceptor.Inject(func(data interface{}) bool {

		if ev, ok := data.(*socket.SessionEvent); ok {

			// Socket各种事件不要往后台发
			switch ev.MsgID {
			case socket.Event_SessionAccepted,
				socket.Event_SessionConnected:
				return true
			}

			// 封包上打上编号
			ev.ClientID = ev.Ses.ID()

			// TODO 按照封包和逻辑固定分发
			BackendAcceptor.Iterate(func(ses cellnet.Session) bool {

				if DebugMode {
					log.Printf("client->backend, msgid: %d clientid: %d data: %v", ev.MsgID, ev.ClientID, ev.Data)
				}

				ses.(socket.RawSession).RawSend(ev.Packet)

				return true
			})

		}

		return false
	})

	ClientAcceptor.Start(address)

}
