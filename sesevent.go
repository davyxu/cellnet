package cellnet

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

type EventType int

const (
	Event_Unknown EventType = iota
	Event_Connected
	Event_ConnectFailed
	Event_Accepted
	Event_AcceptFailed
	Event_Closed
	Event_Recv
	Event_Send
	Event_Post
)

// 会话事件
type Event struct {
	UID int64

	Type EventType // 事件类型

	MsgID uint32      // 消息ID
	Msg   interface{} // 消息对象
	Data  []byte      // 消息序列化后的数据

	Tag         interface{} // 事件的连接, 一个处理流程后被Reset
	TransmitTag interface{} // 接收过程可以传递到发送过程, 不会被清空

	Ses         Session        // 会话
	SendHandler []EventHandler // 发送handler override

	r Result // 出现错误, 将结束ChainCall
}

func (self *Event) Clone() *Event {
	c := &Event{
		UID:         self.UID,
		Type:        self.Type,
		MsgID:       self.MsgID,
		Msg:         self.Msg,
		Tag:         self.Tag,
		TransmitTag: self.TransmitTag,
		Ses:         self.Ses,
		SendHandler: self.SendHandler,
		Data:        make([]byte, len(self.Data)),
	}

	copy(c.Data, self.Data)

	return c
}

func (self *Event) Result() Result {
	return self.r
}

func (self *Event) SetResult(r Result) {
	self.r = r
}

func (self *Event) IsSystemEvent() bool {
	switch self.Type {
	case Event_Connected,
		Event_ConnectFailed,
		Event_Accepted,
		Event_AcceptFailed,
		Event_Closed:
		return true
	}

	return false
}

// 兼容普通消息发送和rpc消息返回, 推荐
func (self *Event) Send(data interface{}) {

	if self.Ses == nil {
		return
	}

	ev := NewEvent(Event_Send, self.Ses)
	ev.Msg = data
	ev.TransmitTag = self.TransmitTag

	self.Ses.RawSend(self.SendHandler, ev)

}

func (self *Event) PeerName() string {
	if self.Ses == nil {
		return ""
	}

	name := self.Ses.FromPeer().Name()
	if name != "" {
		return name
	}

	return self.Ses.FromPeer().Address()
}

func (self *Event) TypeString() string {
	switch self.Type {
	case Event_Recv:
		return "Event_Recv"
	case Event_Send:
		return "Event_Send"
	case Event_Connected:
		return "Event_Connected"
	case Event_ConnectFailed:
		return "Event_ConnectFailed"
	case Event_Accepted:
		return "Event_Accepted"
	case Event_AcceptFailed:
		return "Event_AcceptFailed"
	case Event_Closed:
		return "Event_Closed"
	}

	return fmt.Sprintf("unknown(%d)", self.Type)
}

func (self *Event) SessionID() int64 {
	if self.Ses == nil {
		return 0
	}

	return self.Ses.ID()
}

func (self *Event) MsgSize() int {
	if self.Data == nil {
		return 0
	}

	return len(self.Data)
}

func (self *Event) MsgString() string {

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

func (self *Event) MsgName() string {

	meta := MessageMetaByID(self.MsgID)
	if meta == nil {
		return ""
	}

	return meta.Name
}

func (self *Event) String() string {
	return fmt.Sprintf("#%s(%s) sid: %d MsgID: %d %s | %s Raw: (%d)%v", self.TypeString(), self.PeerName(), self.Ses.ID(), self.MsgID, self.MsgName(), self.MsgString(), self.MsgSize(), self.Data)
}

func (self *Event) FromMessage(msg interface{}) *Event {

	var err error
	self.Data, self.MsgID, err = EncodeMessage(msg)

	if err != nil {
		log.Debugln(err, self.String())
	}

	return self
}

func (self *Event) FromMeta(meta *MessageMeta) *Event {

	if meta != nil {
		self.MsgID = meta.ID
	}

	return self
}

// 根据消息内容, 自动填充其他部分, 以方便输出日志
func (self *Event) Parse() {

	if self.Msg == nil && self.Data != nil && self.MsgID != 0 {

		self.Msg, _ = DecodeMessage(self.MsgID, self.Data)

	} else if self.MsgID == 0 && self.Msg != nil {
		meta := MessageMetaByType(reflect.TypeOf(self.Msg))
		if meta != nil {
			self.MsgID = meta.ID
		}
	}
}

func NewEvent(t EventType, s Session) *Event {
	self := &Event{
		Type: t,
		Ses:  s,
	}

	if EnableHandlerLog {
		self.UID = genSesEvUID()
	}

	return self
}

var evuid int64

func genSesEvUID() int64 {
	atomic.AddInt64(&evuid, 1)
	return evuid
}
