package httpform

import (
	"github.com/davyxu/cellnet/codec"
	"net/http"
)

type httpFormCodec struct {
}

const defaultMemory = 32 * 1024 * 1024

func (self *httpFormCodec) Name() string {
	return "httpform"
}

func (self *httpFormCodec) Encode(msgObj interface{}) (data interface{}, err error) {

	return nil, nil
}

func (self *httpFormCodec) Decode(data interface{}, msgObj interface{}) error {

	req := data.(*http.Request)

	if err := req.ParseForm(); err != nil {
		return err
	}
	req.ParseMultipartForm(defaultMemory)
	if err := mapForm(msgObj, req.Form); err != nil {
		return err
	}

	return nil
}

func init() {

	codec.RegisterCodec(new(httpFormCodec))
}
