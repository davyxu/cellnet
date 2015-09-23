package socket

import (
	"github.com/davyxu/cellnet"
	"log"
	"reflect"
)

// 将PB消息解析封装到闭包中
func RegisterMessage(evq *cellnet.EvQueue, msgIns interface{}, userHandler func(cellnet.Session, interface{})) {

	msgType := reflect.TypeOf(msgIns)

	msgName := msgType.String()

	msgID := cellnet.Name2ID(msgName)

	// 将消息注册到mapper中, 提供反射用
	addMapper(msgName, msgID)

	evq.RegisterCallback(msgID, func(data interface{}) {

		if dv, ok := data.(*DataEvent); ok {

			rawMsg, err := cellnet.ParsePacket(dv.Packet, msgType)

			if err != nil {
				log.Printf("[cellnet] unmarshaling error:\n", err)
				return
			}

			userHandler(dv.Ses, rawMsg)

		}

	})
}
