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

func FetchCodec(name string) Codec {

	return codecByName[name]
}

var (
	ErrMessageNotFound = errors.New("msg not exists")
	ErrCodecNotFound   = errors.New("codec not found")
)

func EncodeMessage(msg interface{}) (data []byte, msgid uint32, err error) {

	// 获取消息元信息
	meta := MessageMetaByType(reflect.TypeOf(msg))
	if meta != nil {
		msgid = meta.ID
	} else {
		return nil, 0, ErrMessageNotFound
	}

	// 将消息编码为字节数组
	data, err = meta.Codec.Encode(msg)

	return data, msgid, err
}

func DecodeMessage(msgid uint32, data []byte) (interface{}, error) {

	// 获取消息元信息
	meta := MessageMetaByID(msgid)

	// 消息没有注册
	if meta == nil {
		return nil, ErrMessageNotFound
	}

	// 创建消息
	msg := reflect.New(meta.Type).Interface()

	// 从字节数组转换为消息
	err := meta.Codec.Decode(data, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
