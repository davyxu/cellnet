package cellnet

import "fmt"

type MatchMsgIDHandler struct {
	msgid uint32
}

func (self *MatchMsgIDHandler) String() string {
	return fmt.Sprintf("MatchMsgIDHandler(msgid:%d %s)", self.msgid, MessageNameByID(self.msgid))
}

func (self *MatchMsgIDHandler) Call(ev *Event) {

	if ev.MsgID != self.msgid {
		ev.SetResult(Result_NextChain)
	}

}

func NewMatchMsgIDHandler(msgid uint32) EventHandler {
	return &MatchMsgIDHandler{
		msgid: msgid,
	}
}
