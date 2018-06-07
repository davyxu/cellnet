package codec

import (
	"github.com/davyxu/cellnet"
)

// 编码消息, 在使用了带内存池的codec中，可以传入session或peer的ContextSet，保存内存池上下文，默认ctx传nil
func EncodeMessage(msg interface{}, ctx cellnet.ContextSet) (data []byte, meta *cellnet.MessageMeta, err error) {

	// 获取消息元信息
	meta = cellnet.MessageMetaByMsg(msg)
	if meta == nil {
		return nil, nil, cellnet.NewErrorContext("msg not exists", msg)
	}

	// 将消息编码为字节数组
	var raw interface{}
	raw, err = meta.Codec.Encode(msg, ctx)

	if err != nil {
		return
	}

	data = raw.([]byte)

	return
}

// 解码消息
func DecodeMessage(msgid int, data []byte) (interface{}, *cellnet.MessageMeta, error) {

	// 获取消息元信息
	meta := cellnet.MessageMetaByID(msgid)

	// 消息没有注册
	if meta == nil {
		return nil, nil, cellnet.NewErrorContext("msg not exists", msgid)
	}

	// 创建消息
	msg := meta.NewType()

	// 从字节数组转换为消息
	err := meta.Codec.Decode(data, msg)

	if err != nil {
		return nil, meta, err
	}

	return msg, meta, nil
}

// Codec.Encode内分配的资源，在必要时可以回收，例如内存池对象
type CodecRecycler interface {
	Free(data interface{}, ctx cellnet.ContextSet)
}

func FreeCodecResource(codec cellnet.Codec, data interface{}, ctx cellnet.ContextSet) {

	if codec == nil {
		return
	}

	if recycler, ok := codec.(CodecRecycler); ok {
		recycler.Free(data, ctx)
	}
}
