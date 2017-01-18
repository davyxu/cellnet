package socket

import "github.com/davyxu/cellnet"

func MessageRegistedCount(evd cellnet.EventDispatcher, msgName string) int {

	msgMeta := cellnet.MessageMetaByName(msgName)
	if msgMeta == nil {
		return 0
	}

	return evd.CountByID(msgMeta.ID)
}

type RegisterMessageContext struct {
	*cellnet.MessageMeta
	*cellnet.CallbackContext
}

// 注册连接消息
func RegisterMessage(evd cellnet.EventDispatcher, msgName string, userHandler func(interface{}, cellnet.Session)) *RegisterMessageContext {

	msgMeta := cellnet.MessageMetaByName(msgName)

	if msgMeta == nil {
		log.Errorf("message register failed, %s", msgName)
		return nil
	}

	ctx := evd.AddCallback(msgMeta.ID, func(data interface{}) {

		if ev, ok := data.(*SessionEvent); ok {

			rawMsg, err := cellnet.ParsePacket(ev.Packet, msgMeta.Type)

			if err != nil {
				log.Errorf("unmarshaling error: %v, raw: %v", err, ev.Packet)
				return
			}

			userHandler(rawMsg, ev.Ses)

		}

	})

	return &RegisterMessageContext{MessageMeta: msgMeta, CallbackContext: ctx}
}
