package json

import (
	"encoding/json"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"net/url"
)

type jsonCodec struct {
}

func (self *jsonCodec) Name() string {
	return "wsjson"
}

func (self *jsonCodec) MimeType() string {
	return "application/json"
}

// 将结构体编码为JSON的字节数组
func (self *jsonCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {

	raw, err := json.Marshal(msgObj)
	if err != nil {
		return nil, err
	}

	data = []byte(url.PathEscape(string(raw)))

	return
}

// 将JSON的字节数组解码为结构体
func (self *jsonCodec) Decode(data interface{}, msgObj interface{}) error {

	decodedStr, err := url.PathUnescape(string(data.([]byte)))
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(decodedStr), msgObj)
}

func init() {

	// 注册编码器
	codec.RegisterCodec(new(jsonCodec))
}
