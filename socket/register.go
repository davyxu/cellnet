package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/log"
	"reflect"
)

// 注册连接消息
func RegisterSessionMessage(eq cellnet.EventQueue, msgIns interface{}, userHandler func(interface{}, cellnet.Session)) *cellnet.MessageMeta {

	msgMeta := cellnet.NewMessageMeta(msgIns)

	// 将消息注册到mapper中, 提供反射用
	MapNameID(msgMeta.Name, msgMeta.ID)

	eq.RegisterCallback(msgMeta.ID, func(data interface{}) {

		if ev, ok := data.(*SessionEvent); ok {

			rawMsg, err := cellnet.ParsePacket(ev.Packet, msgMeta.Type)

			if err != nil {
				log.Errorln("[cellnet] unmarshaling error:\n", err)
				return
			}

			userHandler(rawMsg, ev.Ses)

		}

	})

	return msgMeta
}

// 注册连接消息
func RegisterPeerMessage(eq cellnet.EventQueue, msgIns interface{}, userHandler func(interface{}, cellnet.Peer)) *cellnet.MessageMeta {

	msgMeta := cellnet.NewMessageMeta(msgIns)

	// 将消息注册到mapper中, 提供反射用
	MapNameID(msgMeta.Name, msgMeta.ID)

	eq.RegisterCallback(msgMeta.ID, func(data interface{}) {

		if ev, ok := data.(*PeerEvent); ok {

			rawMsg := reflect.New(msgMeta.Type).Interface()

			userHandler(rawMsg, ev.P)

		}

	})

	return msgMeta
}
