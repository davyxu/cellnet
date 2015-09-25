package gate

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/golang/protobuf/proto"
	"log"
	"reflect"
)

// 连接到Gate的连接器

var gateConnArray []cellnet.Peer

func StartGateConnector(pipe *cellnet.EvPipe, addressList []string) {

	gateConnArray = make([]cellnet.Peer, len(addressList))

	for index, addr := range addressList {

		conn := socket.NewConnector(pipe).Start(addr)
		conn.SetRelayMode(true)
		gateConnArray[index] = conn

		gateIndex := new(int)
		*gateIndex = index

		if DebugMode {
			conn.Inject(func(data interface{}) bool {

				if ev, ok := data.(*socket.SessionEvent); ok {

					// Socket各种事件过滤掉
					switch ev.MsgID {
					case socket.Event_SessionConnected:
						return true
					}

					if ev, ok := data.(*socket.SessionEvent); ok {
						log.Printf("gate->backend, gateindex: %d msgid: %d clientid: %d data: %v", *gateIndex, ev.MsgID, ev.ClientID, ev.Data)
					}
				}

				return true
			})
		}

	}

}

// 注册连接消息
func RegisterSessionMessage(msgIns interface{}, userHandler func(cellnet.Session, int64, interface{})) {

	msgType := reflect.TypeOf(msgIns)

	msgName := msgType.String()

	msgID := cellnet.Name2ID(msgName)

	// 将消息注册到mapper中, 提供反射用
	socket.MapNameID(msgName, msgID)

	for _, conn := range gateConnArray {

		conn.RegisterCallback(msgID, func(data interface{}) {

			if ev, ok := data.(*socket.SessionEvent); ok {

				rawMsg, err := cellnet.ParsePacket(ev.Packet, msgType)

				if err != nil {
					log.Printf("[cellnet] unmarshaling error:\n", err)
					return
				}

				userHandler(ev.Ses, ev.ClientID, rawMsg)

			}

		})
	}

}

// 将消息发送到客户端
func SendToClient(gateSes cellnet.Session, clientid int64, data interface{}) {

	if gateSes == nil {
		return
	}

	gateSes.(socket.RawSession).RelaySend(data, clientid)
}

// 通知网关关闭客户端连接
func CloseClient(gateSes cellnet.Session, clientid int64) {

	if gateSes == nil {
		return
	}

	// 通知关闭
	gateSes.Send(&coredef.CloseClientACK{
		ClientID: proto.Int64(clientid),
	})
}

// 广播所有的客户端
func CloseAllClient() {

	ack := &coredef.CloseClientACK{}

	for _, conn := range gateConnArray {
		ses := conn.(connSesManager).DefaultSession()
		if ses == nil {
			continue
		}

		ses.Send(ack)
	}
}

type connSesManager interface {
	DefaultSession() cellnet.Session
}

// 发送给所有gate的所有客户端
func BroadcastToClient(data interface{}) {

	pkt := cellnet.BuildPacket(data)

	ack := &coredef.BroardcastACK{
		Data:  pkt.Data,
		MsgID: proto.Uint32(pkt.MsgID),
	}

	for _, conn := range gateConnArray {
		ses := conn.(connSesManager).DefaultSession()
		if ses == nil {
			continue
		}

		ses.Send(ack)
	}
}

// 客户端列表
type ClientList map[cellnet.Session][]int64

func (self ClientList) Add(gateSes cellnet.Session, clientid int64) {

	// 事件
	list, ok := self[gateSes]

	// 新建
	if !ok {

		list = make([]int64, 0)

	}

	list = append(list, clientid)
}

func NewClientList() ClientList {
	return make(map[cellnet.Session][]int64)
}

// 发送给指定客户端列表的客户端
func BroadcastToClientList(data interface{}, list ClientList) {

	pkt := cellnet.BuildPacket(data)

	for ses, clientlist := range list {

		ack := &coredef.BroardcastACK{
			Data:  pkt.Data,
			MsgID: proto.Uint32(pkt.MsgID),
		}

		ack.ClientID = clientlist

		ses.Send(ack)
	}

}
