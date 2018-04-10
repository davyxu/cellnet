package http

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

type httpConnector struct {
	peer.CorePeerProperty
	peer.CoreProcBundle
	peer.CoreContextSet
}

func (self *httpConnector) Start() cellnet.Peer {

	return self
}

func (self *httpConnector) Stop() {

}

func (self *httpConnector) Request(method string, raw interface{}) (interface{}, error) {

	// 获取消息元信息
	meta := cellnet.HttpMetaByRequestType(method, reflect.TypeOf(raw))
	if meta == nil {
		return nil, cellnet.NewErrorContext("msg not found", raw)
	}

	// 将消息编码为字节数组
	data, err := meta.RequestCodec.Encode(raw)

	log.Debugf("#http.send(%s) '%s' %s | Message(%s) %s",
		self.Name(),
		meta.Method,
		meta.URL,
		meta.RequestTypeName(),
		cellnet.MessageToString(raw))

	url := fmt.Sprintf("http://%s%s", self.Address(), meta.URL)

	req, err := http.NewRequest(meta.Method, url, data.(io.Reader))

	if err != nil {
		return nil, err
	}

	mimeType := meta.RequestCodec.(interface {
		MimeType() string
	}).MimeType()

	req.Header.Set("Content-Type", mimeType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	msg := reflect.New(meta.ResponseType).Interface()

	err = meta.ResponseCodec.Decode(body, msg)

	log.Debugf("#http.recv(%s) '%s' %s | [%d] Message(%s) %s",
		self.Name(),
		resp.Request.Method,
		meta.URL,
		resp.StatusCode,
		meta.ResponseTypeName(),
		cellnet.MessageToString(msg))

	return msg, err
}

func (self *httpConnector) TypeName() string {
	return "http.Connector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &httpConnector{}

		return p
	})
}
