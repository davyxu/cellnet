package json

import (
	"encoding/json"
	"github.com/davyxu/cellnet"
)

type jsonCodec struct {
}

// 编码器的名称
func (self *jsonCodec) Name() string {
	return "json"
}

// 将结构体编码为JSON的字节数组
func (self *jsonCodec) Encode(msgObj interface{}) ([]byte, error) {

	return json.Marshal(msgObj)

}

// 将JSON的字节数组解码为结构体
func (self *jsonCodec) Decode(data []byte, msgObj interface{}) error {

	return json.Unmarshal(data, msgObj)
}

func init() {

	// 注册编码器
	cellnet.RegisterCodec(new(jsonCodec))
}
