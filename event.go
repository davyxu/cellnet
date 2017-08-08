package cellnet

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

type EventType int32

const (
	Event_None EventType = iota
	Event_Connected
	Event_ConnectFailed
	Event_Accepted
	Event_AcceptFailed
	Event_Closed
	Event_Recv
	Event_Send
)

func (self EventType) String() string {
	switch self {
	case Event_Recv:
		return "recv"
	case Event_Send:
		return "send"
	case Event_Connected:
		return "connected"
	case Event_ConnectFailed:
		return "connectfailed"
	case Event_Accepted:
		return "accepted"
	case Event_AcceptFailed:
		return "acceptfailed"
	case Event_Closed:
		return "closed"
	}

	return fmt.Sprintf("unknown(%d)", self)
}

type Result int32

const (
	Result_OK            Result = iota
	Result_SocketError          // 网络错误
	Result_SocketTimeout        // Socket超时
	Result_PackageCrack         // 封包破损
	Result_CodecError
	Result_RequestClose // 请求关闭
	Result_NextChain
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

	Ses       Session       // 会话
	ChainSend *HandlerChain // 发送handler override

	r Result // 出现错误, 将结束ChainCall

	chainid int64 // 所在链, 调试用
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
		ChainSend:   self.ChainSend,
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

// 兼容普通消息发送和rpc消息返回, 推荐
func (self *Event) Send(data interface{}) {

	if self.Ses == nil {
		return
	}

	ev := NewEvent(Event_Send, self.Ses)
	ev.Msg = data
	ev.TransmitTag = self.TransmitTag

	if self.ChainSend != nil {
		// 由接收方提供的发送链继续传递
		ev.ChainSend = self.ChainSend
	} else {
		// 默认没有, 使用peer的发送链
		ev.ChainSend = self.Ses.FromPeer().ChainSend()
	}

	self.Ses.RawSend(ev)

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

func (self *Event) FromMessage(msg interface{}) *Event {

	var err error
	self.Data, self.MsgID, err = EncodeMessage(msg)

	if err != nil {
		log.Debugln(err, *self)
	}

	// Data+ID / Msg 二选一
	self.Msg = nil

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
