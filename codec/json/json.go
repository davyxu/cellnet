package json

import (
	"encoding/json"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

type jsonCodec struct {
}

// 编码器的名称
func (self *jsonCodec) Name() string {
	return "json"
}

func (self *jsonCodec) MimeType() string {
	return "application/json"
}

// 将结构体编码为JSON的字节数组
func (self *jsonCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {

	return json.Marshal(msgObj)

}

// 将JSON的字节数组解码为结构体
func (self *jsonCodec) Decode(data interface{}, msgObj interface{}) error {

	return json.Unmarshal(data.([]byte), msgObj)
}

func init() {

	// 注册编码器
	codec.RegisterCodec(new(jsonCodec))
}
