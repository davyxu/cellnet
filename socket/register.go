package socket

import "github.com/davyxu/cellnet"

// 注册连接消息
func RegisterMessage(evd cellnet.EventDispatcher, msgName string, userHandler func(interface{}, cellnet.Session)) *cellnet.MessageMeta {

	msgMeta := cellnet.MessageMetaByName(msgName)

	if msgMeta == nil {
		log.Errorf("message register failed, %s", msgName)
		return nil
	}

	evd.RegisterCallback(msgMeta.ID, func(data interface{}) {

		if ev, ok := data.(*SessionEvent); ok {

			rawMsg, err := cellnet.ParsePacket(ev.Packet, msgMeta.Type)

			if err != nil {
				log.Errorln("unmarshaling error:\n", err)
				return
			}

			userHandler(rawMsg, ev.Ses)

		}

	})

	return msgMeta
}
