package cellnet

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

type EventType int

const (
	SessionEvent_Unknown EventType = iota
	SessionEvent_Connected
	SessionEvent_ConnectFailed
	SessionEvent_Accepted
	SessionEvent_AcceptFailed
	SessionEvent_Recv
	SessionEvent_Send
)

// 会话事件
type SessionEvent struct {
	Type  EventType    // 事件类型
	MsgID uint32       // 消息ID
	Msg   interface{}  // 消息对象
	Meta  *MessageMeta // 消息扩展内容
	Data  []byte       // 消息序列化后的数据
	Ses   Session      // 会话
	Tag   interface{}  // 事件的连接
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

func (self *SessionEvent) DirString() string {
	switch self.Type {
	case SessionEvent_Recv:
		return "recv"
	case SessionEvent_Send:
		return "send"
	case SessionEvent_Connected:
		return "connected"
	case SessionEvent_Accepted:
		return "accepted"
	}

	return "unknown"
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
	}

	return "unknown"
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

	return self.Msg.(proto.Message).String()
}

func (self *SessionEvent) MsgName() string {

	if self.Meta == nil {
		return ""
	}

	return self.Meta.Name
}

func (self *SessionEvent) String() string {
	return fmt.Sprintf("#%s(%s) MsgID: %d %s | %s Raw: (%d)%v", self.TypeString(), self.PeerName(), self.MsgID, self.MsgName(), self.MsgString(), self.MsgSize(), self.Data)
}

func (self *SessionEvent) FromMessage(msg interface{}) *SessionEvent {
	self.Data, self.Meta = BuildPacket(msg)
	if self.Meta != nil {
		self.MsgID = self.Meta.ID
	}

	return self
}

func (self *SessionEvent) FromMeta(meta *MessageMeta) *SessionEvent {
	self.Meta = meta
	if meta != nil {
		self.MsgID = self.Meta.ID
	}

	return self
}

func NewSessionEvent(t EventType, s Session) *SessionEvent {
	return &SessionEvent{
		Type: t,
		Ses:  s,
	}
}
