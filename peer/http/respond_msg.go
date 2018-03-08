package http

import (
	"errors"
	"github.com/davyxu/cellnet"
	"net/http"
	"reflect"
)

type Message struct {
	StatusCode int
	Msg        interface{}
}

func (self *Message) WriteRespond(ses *httpSession) error {
	peerInfo := ses.Peer().(cellnet.PeerProperty)

	msg := self.Msg

	log.Debugf("#recv(%s) http.%s %s | [%d] Message(%s) %s",
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

	// 将消息编码为字节数组
	var data interface{}
	data, err := meta.ResponseCodec.Encode(msg)

	if err != nil {
		return err
	}

	ses.resp.WriteHeader(http.StatusOK)
	ses.resp.Write(data.([]byte))

	return nil
}
