package cellnet

import (
	xnet "github.com/davyxu/x/net"
	"net"
)

// 长连接
type Session interface {

	// 获得原始的Socket连接
	Raw() interface{}

	// 获得Session归属的Peer
	Peer() Peer

	// 发送消息，消息需要以指针格式传入
	Send(msg interface{})

	// 断开
	Close()

	// 标示ID
	ID() int64
}

// 直接发送数据时，将*RawPacket作为Send参数
type RawPacket struct {
	MsgData []byte
	MsgID   int
	MsgName string
}

func (self *RawPacket) Message() interface{} {

	// 获取消息元信息
	meta := MessageMetaByID(self.MsgID)

	// 消息没有注册
	if meta == nil {
		return struct{}{}
	}

	// 创建消息
	msg := meta.NewType()

	// 从字节数组转换为消息
	err := meta.Codec.Decode(self.MsgData, msg)
	if err != nil {
		return struct{}{}
	}

	return msg
}

// 修复ws没有实现所有net.conn方法，导致无法获取客服端地址问题.
type RemoteAddr interface {
	RemoteAddr() net.Addr
}

// 获取session远程的地址
func GetRemoteAddrss(ses Session) string {

	if ses == nil {
		return ""
	}

	if c, ok := ses.Raw().(RemoteAddr); ok {
		return c.RemoteAddr().String()
	}

	return ""
}

func GetRemoteHost(ses Session) string {
	addr := GetRemoteAddrss(ses)
	host, _, err := xnet.SpliteAddress(addr)
	if err == nil {
		return host
	}

	return ""
}
