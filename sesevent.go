package cellnet

import (
	"fmt"
	"reflect"
)

type EventType int

const (
	SessionEvent_Unknown EventType = iota
	SessionEvent_Connected
	SessionEvent_ConnectFailed
	SessionEvent_Accepted
	SessionEvent_AcceptFailed
	SessionEvent_Closed
	SessionEvent_Recv
	SessionEvent_Send
)

// 会话事件
type SessionEvent struct {
	Type EventType // 事件类型

	MsgID uint32      // 消息ID
	Msg   interface{} // 消息对象
	Data  []byte      // 消息序列化后的数据

	Tag interface{} // 事件的连接

	Ses         Session      // 会话
	SendHandler EventHandler // 发送handler override

	EndRecvLoop bool // 停止消息接收循环
}

func (self *SessionEvent) IsSystemEvent() bool {
	switch self.Type {
	case SessionEvent_Connected,
		SessionEvent_ConnectFailed,
		SessionEvent_Accepted,
		SessionEvent_AcceptFailed,
		SessionEvent_Closed:
		return true
	}

	return false
}

// 兼容普通消息发送和rpc消息返回, 推荐
func (self *SessionEvent) Send(data interface{}) {

	if self.Ses == nil {
		return
	}

	self.Reset(SessionEvent_Send)
	self.Msg = data

	self.Ses.RawSend(self.SendHandler, self)

}

func (self *SessionEvent) Reset(t EventType) {
	self.Type = t
	self.MsgID = 0
	self.Msg = nil
	self.Data = nil
	self.Tag = nil
}

func (self *SessionEvent) PeerName() string {
	if self.Ses == nil {
		return ""
	}

	name := self.Ses.FromPeer().Name()
	if name != "" {
		return name
	}

	return self.Ses.FromPeer().Address()
}

func (self *SessionEvent) TypeString() string {
	switch self.Type {
	case SessionEvent_Recv:
		return "SessionEvent_Recv"
	case SessionEvent_Send:
		return "SessionEvent_Send"
	case SessionEvent_Connected:
		return "SessionEvent_Connected"
	case SessionEvent_ConnectFailed:
		return "SessionEvent_ConnectFailed"
	case SessionEvent_Accepted:
		return "SessionEvent_Accepted"
	case SessionEvent_Closed:
		return "SessionEvent_Closed"
	}

	return fmt.Sprintf("unknown(%d)", self.Type)
}

func (self *SessionEvent) SessionID() int64 {
	if self.Ses == nil {
		return 0
	}

	return self.Ses.ID()
}

func (self *SessionEvent) MsgSize() int {
	if self.Data == nil {
		return 0
	}

	return len(self.Data)
}

func (self *SessionEvent) MsgString() string {

	if self.Msg == nil {
		return ""
	}

	if stringer, ok := self.Msg.(interface {
		String() string
	}); ok {
		return stringer.String()
	}

	return ""
}

func (self *SessionEvent) MsgName() string {

	meta := MessageMetaByID(self.MsgID)
	if meta == nil {
		return ""
	}

	return meta.Name
}

func (self *SessionEvent) String() string {
	return fmt.Sprintf("#%s(%s) MsgID: %d %s | %s Raw: (%d)%v", self.TypeString(), self.PeerName(), self.MsgID, self.MsgName(), self.MsgString(), self.MsgSize(), self.Data)
}

func (self *SessionEvent) FromMessage(msg interface{}) *SessionEvent {

	meta := MessageMetaByName(MessageFullName(reflect.TypeOf(msg)))
	if meta != nil {
		self.MsgID = meta.ID
	}

	if meta.Codec == nil {
		log.Errorf("message codec not found: %s", meta.Name)
		return self
	}

	var err error
	self.Data, err = meta.Codec.Encode(msg)

	if err != nil {
		log.Errorln(err)
	}

	return self
}

func (self *SessionEvent) FromMeta(meta *MessageMeta) *SessionEvent {

	if meta != nil {
		self.MsgID = meta.ID
	}

	return self
}

func NewSessionEvent(t EventType, s Session) *SessionEvent {
	return &SessionEvent{
		Type: t,
		Ses:  s,
	}
}
