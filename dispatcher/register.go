package dispatcher

import (
	"github.com/davyxu/cellnet"
	"log"
	"reflect"
)

// 将PB消息解析封装到闭包中
func RegisterMessage(disp *DataDispatcher, msgIns interface{}, userHandler func(cellnet.CellID, interface{})) {

	msgType := reflect.TypeOf(msgIns)

	msgName := msgType.String()

	msgID := cellnet.Name2ID(msgName)

	// 将消息注册到mapper中, 提供反射用
	addMapper(msgName, msgID)

	log.Printf("[dispatcher] #regmsg %s(%d)", msgName, msgID)

	disp.RegisterCallback(msgID, func(data interface{}) {

		switch sd := data.(type) {
		case cellnet.SessionPacket:

			rawMsg, err := cellnet.ParsePacket(sd.GetPacket(), msgType)

			if err != nil {
				log.Printf("[cellnet] unmarshaling error:\n", err)
				return
			}

			userHandler(sd.GetSession(), rawMsg)
		case cellnet.SessionEvent:
			userHandler(sd.GetSession(), nil)
		}

	})
}

// 将PB消息解析封装到闭包中
func RegisterRemoteCall(disp *DataDispatcher, msgIns interface{}, userHandler func(interface{}, cellnet.RPCResponse)) {

	msgType := reflect.TypeOf(msgIns)

	msgName := msgType.String()

	msgID := cellnet.Name2ID(msgName)

	// 将消息注册到mapper中, 提供反射用
	addMapper(msgName, msgID)

	//log.Printf("[dispatcher] #regmsg %s(%d)", msgName, msgID)

	disp.RegisterCallback(msgID, func(data interface{}) {

		switch sd := data.(type) {
		case cellnet.RPCResponse:

			rawMsg, err := cellnet.ParsePacket(sd.GetPacket(), msgType)

			if err != nil {
				log.Printf("[cellnet] unmarshaling error:\n", err)
				return
			}

			userHandler(rawMsg, sd)

		case cellnet.SessionEvent:
			userHandler(sd.GetSession(), nil)
		}

	})
}
