package json

import (
	"encoding/json"
	cellcodec "github.com/davyxu/cellnet/codec"
	"github.com/davyxu/x/container"
)

type jsonCodec struct {
}

// 编码器的名称
func (self *jsonCodec) Name() string {
	return "json"
}

// 将结构体编码为JSON的字节数组
func (self *jsonCodec) Encode(msgObj any, ps *xcontainer.Mapper) (data any, err error) {

	return json.Marshal(msgObj)

}

// 将JSON的字节数组解码为结构体
func (self *jsonCodec) Decode(data, msgObj any) error {

	return json.Unmarshal(data.([]byte), msgObj)
}

func init() {

	// 注册编码器
	cellcodec.Register(new(jsonCodec))
}
