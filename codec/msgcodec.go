package codec

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

// 编码消息
func EncodeMessage(msg interface{}) (data []byte, meta *cellnet.MessageMeta, err error) {

	// 获取消息元信息
	meta = cellnet.MessageMetaByType(reflect.TypeOf(msg))
	if meta == nil {
		return nil, nil, cellnet.NewErrorContext("msg not exists", msg)
	}

	// 将消息编码为字节数组
	var raw interface{}
	raw, err = meta.Codec.Encode(msg)

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
