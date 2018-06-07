package httpjson

import (
	"bytes"
	"encoding/json"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"io"
	"io/ioutil"
	"net/http"
)

type httpjsonCodec struct {
}

// 编码器的名称
func (self *httpjsonCodec) Name() string {
	return "httpjson"
}

func (self *httpjsonCodec) MimeType() string {
	return "application/json"
}

// 将结构体编码为JSON的字节数组
func (self *httpjsonCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {

	bdata, err := json.Marshal(msgObj)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bdata), nil
}

// 将JSON的字节数组解码为结构体
func (self *httpjsonCodec) Decode(data interface{}, msgObj interface{}) error {

	var reader io.Reader
	switch v := data.(type) {
	case *http.Request:
		reader = v.Body
	case io.Reader:
		reader = v
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, msgObj)
}

func init() {

	// 注册编码器
	codec.RegisterCodec(new(httpjsonCodec))
}
