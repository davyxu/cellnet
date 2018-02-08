package http

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"net/http"
	"reflect"
)

type StatusCode int

type HttpContext interface {
	Request() *http.Request
	Response() http.ResponseWriter
}

type MessageProc struct {
}

var (
	errNotHandled = errors.New("Request not handled")
	errNotFound   = errors.New("404 Not found")
)

func (MessageProc) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	httpContext := ses.(HttpContext)
	req := httpContext.Request()

	meta := cellnet.HttpMetaByMethodURL(req.Method, req.URL.Path)
	if meta != nil {

		msg = reflect.New(meta.RequestType).Interface()

		if err = meta.RequestCodec.Decode(req, msg); err != nil {
			return
		}

		return

	}

	return nil, errNotHandled
}

func (MessageProc) OnSendMessage(ses cellnet.Session, raw interface{}) error {

	httpContext := ses.(HttpContext)
	resp := httpContext.Response()

	switch msg := raw.(type) {
	case StatusCode:
		resp.WriteHeader(int(msg))
	default:

		// 获取消息元信息
		meta := cellnet.HttpMetaByResponseType(httpContext.Request().Method, reflect.TypeOf(msg))
		if meta == nil {
			return codec.ErrMessageNotFound
		}

		// 将消息编码为字节数组
		var data interface{}
		data, err := meta.ResponseCodec.Encode(msg)

		if err != nil {
			return err
		}

		resp.Write(data.([]byte))
	}

	return nil
}
