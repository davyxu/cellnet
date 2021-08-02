package json

import (
	"encoding/json"
	cellcodec "github.com/davyxu/cellnet/codec"
	xframe "github.com/davyxu/x/frame"
)

type jsonCodec struct {
}

// 编码器的名称
func (self *jsonCodec) Name() string {
	return "json"
}

// 将结构体编码为JSON的字节数组
func (self *jsonCodec) Encode(msgObj interface{}, ps *xframe.Mapper) (data interface{}, err error) {

	return json.Marshal(msgObj)

}

// 将JSON的字节数组解码为结构体
func (self *jsonCodec) Decode(data interface{}, msgObj interface{}) error {

	return json.Unmarshal(data.([]byte), msgObj)
}

func init() {

	// 注册编码器
	cellcodec.Register(new(jsonCodec))
}
