package cellnet

type MatchMsgIDHandler struct {
	msgid uint32

	hlist []EventHandler
}

func (self *MatchMsgIDHandler) Call(ev *SessionEvent) {

	if ev.MsgID == self.msgid {
		HandlerChainCall(self.hlist, ev)
	}

}

func NewMatchMsgIDHandler(msgid uint32, hlist ...EventHandler) EventHandler {
	return &MatchMsgIDHandler{
		msgid: msgid,
		hlist: hlist,
	}
}
