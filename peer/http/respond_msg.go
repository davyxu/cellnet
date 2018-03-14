package http

import (
	"errors"
	"github.com/davyxu/cellnet"
	"net/http"
	"reflect"
)

type MessageRespond struct {
	StatusCode int
	Msg        interface{}
}

func (self *MessageRespond) WriteRespond(ses *httpSession) error {
	peerInfo := ses.Peer().(cellnet.PeerProperty)

	msg := self.Msg

	log.Debugf("#http.recv(%s) '%s' %s | [%d] Message(%s) %s",
		peerInfo.Name(),
		ses.req.Method,
		ses.req.URL.Path,
		self.StatusCode,
		cellnet.MessageToName(msg),
		cellnet.MessageToString(msg))

	// 获取消息元信息
	meta := cellnet.HttpMetaByResponseType(ses.req.Method, reflect.TypeOf(msg))
	if meta == nil {
		return errors.New("message not found:" + reflect.TypeOf(msg).Name())
	}

	if meta.ResponseCodec == nil {
		return errors.New("ResponseCodec not found:" + reflect.TypeOf(msg).Name())
	}

	// 将消息编码为字节数组
	var data interface{}
	data, err := meta.ResponseCodec.Encode(msg)

	if err != nil {
		return err
	}

	ses.resp.Header().Set("Content-Type", meta.ResponseCodec.MimeType()+";charset=UTF-8")
	ses.resp.WriteHeader(http.StatusOK)
	ses.resp.Write(data.([]byte))

	return nil
}
