package cellnet

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
