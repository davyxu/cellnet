package router

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
)

// Backend的各种服务器端使用以下代码

var routerConnArray []cellnet.Peer

type relayEvent struct {
	*socket.SessionEvent

	ClientID int64
}

const defaultReconnectSec = 2

// 后台服务器到router的连接
func StartBackendConnector(pipe cellnet.EventPipe, addressList []string, peerName string, svcName string) {

	routerConnArray = make([]cellnet.Peer, len(addressList))

	if len(addressList) == 0 {
		log.Warnf("empty router address list")
		return
	}

	for index, addr := range addressList {

		peer := socket.NewConnector(pipe)
		peer.SetName(peerName)

		peer.(cellnet.Connector).SetAutoReconnectSec(defaultReconnectSec)

		peer.Start(addr)

		routerConnArray[index] = peer

		// 连上网关时, 发送自己的服务器名字进行注册
		socket.RegisterSessionMessage(peer, "coredef.SessionConnected", func(content interface{}, ses cellnet.Session) {

			ses.Send(&coredef.RegisterRouterBackendACK{
				Name: svcName,
			})

		})

		// 广播
		socket.RegisterSessionMessage(peer, "coredef.UpstreamACK", func(content interface{}, ses cellnet.Session) {
			msg := content.(*coredef.UpstreamACK)

			// 生成派发的消息

			// TODO 用PostData防止多重嵌套?
			// 调用已注册的回调
			peer.CallData(&relayEvent{
				SessionEvent: socket.NewSessionEvent(msg.MsgID, ses, msg.Data),
				ClientID:     msg.ClientID,
			})

		})

	}

}

// 注册连接消息
func RegisterMessage(msgName string, userHandler func(interface{}, cellnet.Session, int64)) {

	msgMeta := cellnet.MessageMetaByName(msgName)

	if msgMeta == nil {
		log.Errorf("message register failed, %s", msgName)
		return
	}

	for _, conn := range routerConnArray {

		conn.RegisterCallback(msgMeta.ID, func(data interface{}) {

			if ev, ok := data.(*relayEvent); ok {

				log.Debugf("router->backend, msg: %s(%d) clientid: %d", getMsgName(ev.MsgID), ev.MsgID, ev.ClientID)

				rawMsg, err := cellnet.ParsePacket(ev.Packet, msgMeta.Type)

				if err != nil {
					log.Errorln("unmarshaling error:\n", err)
					return
				}

				userHandler(rawMsg, ev.Ses, ev.ClientID)

			}

		})
	}

}

// 将消息发送到客户端
func SendToClient(routerSes cellnet.Session, clientid int64, data interface{}) {

	if routerSes == nil {
		return
	}

	userpkt, _ := cellnet.BuildPacket(data)

	log.Debugf("backend->router, msg: %s(%d) clientid: %d", getMsgName(userpkt.MsgID), userpkt.MsgID, clientid)

	routerSes.Send(&coredef.DownstreamACK{
		Data:     userpkt.Data,
		MsgID:    userpkt.MsgID,
		ClientID: []int64{clientid},
	})
}

// 通知网关关闭客户端连接
func CloseClient(routerSes cellnet.Session, clientid int64) {

	if routerSes == nil {
		return
	}

	log.Debugf("backend->router, CloseClient clientid: %d", clientid)

	// 通知关闭
	routerSes.Send(&coredef.CloseClientACK{
		ClientID: clientid,
	})
}

// 广播所有的客户端
func CloseAllClient() {

	ack := &coredef.CloseClientACK{}

	for _, conn := range routerConnArray {
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

// 发送给所有router的所有客户端
func BroadcastToClient(data interface{}) {

	pkt, _ := cellnet.BuildPacket(data)

	ack := &coredef.DownstreamACK{
		Data:  pkt.Data,
		MsgID: pkt.MsgID,
	}

	for _, conn := range routerConnArray {
		ses := conn.(connSesManager).DefaultSession()
		if ses == nil {
			continue
		}

		ses.Send(ack)
	}
}

// 客户端列表
type ClientList map[cellnet.Session][]int64

func (self ClientList) Add(routerSes cellnet.Session, clientid int64) {

	// 事件
	list, ok := self[routerSes]

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
