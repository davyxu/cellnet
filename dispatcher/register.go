package dispatcher

import (
	"github.com/davyxu/cellnet"
	"log"
	"reflect"
)

// 将PB消息解析封装到闭包中
func RegisterMessage(disp *PacketDispatcher, msgIns interface{}, userHandler func(cellnet.CellID, interface{})) {

	msgType := reflect.TypeOf(msgIns)

	msgName := msgType.String()

	msgID := cellnet.Name2ID(msgName)

	disp.RegisterCallback(msgID, func(ses cellnet.CellID, pkt *cellnet.Packet) {

		rawMsg, err := cellnet.ParsePacket(pkt, msgType)

		if err != nil {
			log.Printf("unmarshaling error:\n", err)
			return
		}

		userHandler(ses, rawMsg)
	})
}
