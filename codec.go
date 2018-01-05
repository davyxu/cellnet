package cellnet

import (
	"errors"
	"reflect"
)

type Codec interface {
	// 将数据转换为字节数组
	Encode(interface{}) ([]byte, error)

	// 将字节数组转换为数据
	Decode([]byte, interface{}) error

	// 编码器的名字
	Name() string
}

var codecByName = map[string]Codec{}

func RegisterCodec(c Codec) {

	if _, ok := codecByName[c.Name()]; ok {
		panic("duplicate codec: " + c.Name())
	}

	codecByName[c.Name()] = c
}

func GetCodec(name string) Codec {

	return codecByName[name]
}

func MustGetCodec(name string) Codec {
	codec := GetCodec(name)

	if codec == nil {
		panic("codec not register! " + name)
	}

	return codec
}

var (
	ErrMessageNotFound = errors.New("msg not exists")
)

func EncodeMessage(msg interface{}) (data []byte, meta *MessageMeta, err error) {

	// 获取消息元信息
	meta = MessageMetaByType(reflect.TypeOf(msg))
	if meta == nil {
		return nil, nil, ErrMessageNotFound
	}

	// 将消息编码为字节数组
	data, err = meta.Codec.Encode(msg)

	return
}

func DecodeMessage(msgid int, data []byte) (interface{}, *MessageMeta, error) {

	// 获取消息元信息
	meta := MessageMetaByID(msgid)

	// 消息没有注册
	if meta == nil {
		return nil, nil, ErrMessageNotFound
	}

	// 创建消息
	msg := reflect.New(meta.Type).Interface()

	// 从字节数组转换为消息
	err := meta.Codec.Decode(data, msg)
	if err != nil {
		return nil, meta, err
	}

	return msg, meta, nil
}
