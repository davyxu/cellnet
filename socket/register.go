package socket

import (
	"github.com/davyxu/cellnet"
	"log"
	"reflect"
)

// 注册连接消息
func RegisterSessionMessage(em cellnet.EventManager, msgIns interface{}, userHandler func(cellnet.Session, interface{})) uint32 {

	msgType := reflect.TypeOf(msgIns)

	msgName := msgType.String()

	msgID := cellnet.Name2ID(msgName)

	// 将消息注册到mapper中, 提供反射用
	MapNameID(msgName, msgID)

	em.RegisterCallback(msgID, func(data interface{}) {

		if ev, ok := data.(*SessionEvent); ok {

			rawMsg, err := cellnet.ParsePacket(ev.Packet, msgType)

			if err != nil {
				log.Printf("[cellnet] unmarshaling error:\n", err)
				return
			}

			userHandler(ev.Ses, rawMsg)

		}

	})

	return uint32(msgID)
}

// 注册连接消息
func RegisterPeerMessage(em cellnet.EventManager, msgIns interface{}, userHandler func(cellnet.Peer, interface{})) {

	msgType := reflect.TypeOf(msgIns)

	msgName := msgType.String()

	msgID := cellnet.Name2ID(msgName)

	// 将消息注册到mapper中, 提供反射用
	MapNameID(msgName, msgID)

	em.RegisterCallback(msgID, func(data interface{}) {

		if ev, ok := data.(*PeerEvent); ok {

			rawMsg := reflect.New(msgType).Interface()

			userHandler(ev.P, rawMsg)

		}

	})
}
