package gate

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
)

// 连接到Gate的连接器

var gateConnArray []cellnet.Peer

type relayEvent struct {
	*socket.SessionEvent

	ClientID int64
}

func StartGateConnector(pipe cellnet.EventPipe, addressList []string) {

	gateConnArray = make([]cellnet.Peer, len(addressList))

	for index, addr := range addressList {

		conn := socket.NewConnector(pipe).Start(addr)
		gateConnArray[index] = conn

		gateIndex := new(int)
		*gateIndex = index

		// 广播
		socket.RegisterSessionMessage(conn, "coredef.UpstreamACK", func(content interface{}, ses cellnet.Session) {
			msg := content.(*coredef.UpstreamACK)

			// 生成派发的消息

			// TODO 用PostData防止多重嵌套?
			// 调用已注册的回调
			conn.CallData(&relayEvent{
				SessionEvent: socket.NewSessionEvent(msg.MsgID, ses, msg.Data),
				ClientID:     msg.ClientID,
			})

		})

	}

}

// 注册连接消息
func RegisterSessionMessage(msgName string, userHandler func(interface{}, cellnet.Session, int64)) {

	msgMeta := cellnet.MessageMetaByName(msgName)

	for _, conn := range gateConnArray {

		conn.RegisterCallback(msgMeta.ID, func(data interface{}) {

			if ev, ok := data.(*relayEvent); ok {

				rawMsg, err := cellnet.ParsePacket(ev.Packet, msgMeta.Type)

				if err != nil {
					log.Errorln("[gate] unmarshaling error:\n", err)
					return
				}

				userHandler(rawMsg, ev.Ses, ev.ClientID)

			}

		})
	}

}

// 将消息发送到客户端
func SendToClient(gateSes cellnet.Session, clientid int64, data interface{}) {

	if gateSes == nil {
		return
	}

	userpkt, _ := cellnet.BuildPacket(data)

	gateSes.Send(&coredef.DownstreamACK{
		Data:     userpkt.Data,
		MsgID:    userpkt.MsgID,
		ClientID: []int64{clientid},
	})
}

// 通知网关关闭客户端连接
func CloseClient(gateSes cellnet.Session, clientid int64) {

	if gateSes == nil {
		return
	}

	// 通知关闭
	gateSes.Send(&coredef.CloseClientACK{
		ClientID: clientid,
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

	pkt, _ := cellnet.BuildPacket(data)

	ack := &coredef.DownstreamACK{
		Data:  pkt.Data,
		MsgID: pkt.MsgID,
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

	pkt, _ := cellnet.BuildPacket(data)

	for ses, clientlist := range list {

		ack := &coredef.DownstreamACK{
			Data:  pkt.Data,
			MsgID: pkt.MsgID,
		}

		ack.ClientID = clientlist

		ses.Send(ack)
	}

}
