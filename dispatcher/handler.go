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
func PeerHandler(disp *DataDispatcher) func(cellnet.CellID, interface{}) {

	return func(peer cellnet.CellID, peerev interface{}) {

		switch v := peerev.(type) {
		case ltvsocket.EventCreateSession: // 新的连接生成

			ltvsocket.SpawnSession(v.Stream, v.Type, func(ses cellnet.CellID, sesev interface{}) {

				switch ev := sesev.(type) {

				case ltvsocket.EventNewSession:

					if ev.Type == ltvsocket.SessionAccepted {
						disp.Call(ev.Session, msgAccepted, nil)
					} else {
						disp.Call(ev.Session, msgConnected, nil)
					}

				case ltvsocket.EventClose: // 断开转通知
					disp.Call(ev.Session, msgClosed, nil)
				case ltvsocket.EventData: // 收
					disp.Call(ev.Session, ev.Packet.ID(), ev.Packet)
				}

			})

		case errInterface:
			log.Println(cellnet.ReflectContent(v))
		}

	}
}
