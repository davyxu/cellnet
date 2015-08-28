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

	//log.Printf("[dispatcher] #regmsg %s(%d 0x%x)", msgName, msgID, msgID)

	disp.RegisterCallback(msgID, func(ses cellnet.CellID, data interface{}) {

		if data == nil {

			userHandler(ses, nil)

		} else {
			pkt := data.(*cellnet.Packet)

			rawMsg, err := cellnet.ParsePacket(pkt, msgType)

			if err != nil {
				log.Printf("[cellnet] unmarshaling error:\n", err)
				return
			}

			userHandler(ses, rawMsg)

		}

	})
}
