package dispatcher

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/ltvsocket"
	"github.com/davyxu/cellnet/proto/coredef"
	"log"
)

type errInterface interface {
	Error() string
}

var (
	msgConnected = cellnet.Type2ID(&coredef.ConnectedACK{})
	msgAccepted  = cellnet.Type2ID(&coredef.AcceptedACK{})
	msgClosed    = cellnet.Type2ID(&coredef.ClosedACK{})
)

// 处理Peer的新会话及会话的消息处理
func PeerHandler(disp *DataDispatcher) func(interface{}) {

	return func(peerev interface{}) {

		switch v := peerev.(type) {
		case ltvsocket.SocketCreateSession: // 新的连接生成

			ltvsocket.SpawnSession(v.Stream, v.Type, func(sesev interface{}) {

				switch ev := sesev.(type) {

				case ltvsocket.SocketNewSession:

					if ev.Type == ltvsocket.SessionAccepted {
						disp.Call(msgAccepted, sesev)
					} else {
						disp.Call(msgConnected, sesev)
					}

				case ltvsocket.SocketClose: // 断开转通知
					disp.Call(msgClosed, sesev)
				case cellnet.SessionPacket: // 收

					disp.Call(int(ev.GetPacket().MsgID), sesev)
				}

			})

		case errInterface:
			log.Println(cellnet.ReflectContent(v))
		}

	}
}

func EventHandler(disp *DataDispatcher) func(interface{}) {

	return func(peerev interface{}) {

	}

}
