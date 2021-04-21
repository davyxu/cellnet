package cellcodec

import (
	"fmt"
	"github.com/davyxu/cellnet"
	cellevent "github.com/davyxu/cellnet/event"
	cellmeta "github.com/davyxu/cellnet/meta"
	xframe "github.com/davyxu/x/frame"
)

func init() {
	cellevent.InternalDecodeHandler = func(ev cellnet.Event) (msg interface{}) {
		msg, _, _ = Decode(ev.MessageID(), ev.MessageData())
		return
	}
}

// 编码消息, 在使用了带内存池的codec中，可以传入session或peer的ContextSet，保存内存池上下文，默认ctx传nil
func Encode(msg interface{}, ps *xframe.PropertySet) (data []byte, meta *cellmeta.Meta, err error) {

	// 获取消息元信息
	meta = cellmeta.MetaByMsg(msg)
	if meta == nil {
		return nil, nil, fmt.Errorf("msg not exists: %+v", msg)
	}

	// 将消息编码为字节数组
	var raw interface{}
	raw, err = meta.Codec.Encode(msg, ps)

	if err != nil {
		return
	}

	data = raw.([]byte)

	return
}

// 解码消息
func DecodeByName(msgName string, data []byte) (interface{}, *cellmeta.Meta, error) {

	// 获取消息元信息
	meta := cellmeta.MetaByFullName(msgName)

	// 消息没有注册
	if meta == nil {
		return nil, nil, fmt.Errorf("msg not exists: %s", msgName)
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

// 解码消息
func Decode(msgid int, data []byte) (interface{}, *cellmeta.Meta, error) {

	// 获取消息元信息
	meta := cellmeta.MetaByID(msgid)

	// 消息没有注册
	if meta == nil {
		return nil, nil, fmt.Errorf("msg not exists: %d", msgid)
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
	Free(data interface{}, ps *xframe.PropertySet)
}

func Free(codec cellnet.Codec, data interface{}, ps *xframe.PropertySet) {

	if codec == nil {
		return
	}

	if recycler, ok := codec.(CodecRecycler); ok {
		recycler.Free(data, ps)
	}
}
