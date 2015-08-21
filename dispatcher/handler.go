package dispatcher

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/ltvsocket"
	"github.com/golang/protobuf/proto"
	"log"
)

const (
	EventNewSession = 1
)

type errInterface interface {
	Error() string
}

// 处理Peer的新会话及会话的消息处理
func PeerHandler(disp *PacketDispatcher) func(cellnet.CellID, interface{}) {

	return func(peer cellnet.CellID, peerev interface{}) {

		switch v := peerev.(type) {
		case cellnet.IPacketStream: // 新的连接生成

			ltvsocket.SpawnSession(v, func(ses cellnet.CellID, sesev interface{}) {

				switch data := sesev.(type) {

				case cellnet.EventInit: // 初始化转通知
					disp.Call(ses, &cellnet.Packet{MsgID: EventNewSession})
				case *cellnet.Packet: // 收

					disp.Call(ses, data)

				case proto.Message: // 发
					v.Write(cellnet.BuildPacket(data))
				}

			})

		case errInterface:
			log.Println(cellnet.ReflectContent(v))
		}

	}
}
