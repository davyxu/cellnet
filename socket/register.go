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

type ParsePacketHandler struct {
	meta        *cellnet.MessageMeta
	userHandler func(interface{}, cellnet.Session)
}

func (self *ParsePacketHandler) Call(evid int, data interface{}) error {

	ev := data.(*SessionEvent)

	rawMsg, err := cellnet.ParsePacket(ev.Packet, self.meta.Type)

	if err != nil {
		log.Errorf("unmarshaling error: %v, raw: %v", err, ev.Packet)
		return nil
	}

	self.userHandler(rawMsg, ev.Ses)

	return nil
}

func NewParsePacketHandler(meta *cellnet.MessageMeta, userHandler func(interface{}, cellnet.Session)) cellnet.Handler {
	return &ParsePacketHandler{
		meta:        meta,
		userHandler: userHandler,
	}
}

func RegisterHandler(dh *DispatcherHandler, msgName string, userHandler func(interface{}, cellnet.Session)) {

	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		log.Errorf("message register failed, %s", msgName)
		return
	}

	dh.Add(meta.ID, NewParsePacketHandler(meta, userHandler))

}
