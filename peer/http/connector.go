package http

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/peer"
	"io"
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

func getCodec(codecName string) cellnet.Codec {

	if codecName == "" {
		codecName = "httpjson"
	}

	return codec.MustGetCodec(codecName)
}

func getTypeName(msg interface{}) string {
	if msg == nil {
		return ""
	}

	return reflect.TypeOf(msg).Elem().Name()
}

func (self *httpConnector) Request(method, path string, param *cellnet.HTTPRequest) error {

	// 将消息编码为字节数组
	reqCodec := getCodec(param.REQCodecName)
	data, err := reqCodec.Encode(param.REQMsg, nil)

	if log.IsDebugEnabled() {
		log.Debugf("#http.send(%s) '%s' %s | Message(%s) %s",
			self.Name(),
			method,
			path,
			getTypeName(param.REQMsg),
			cellnet.MessageToString(param.REQMsg))
	}

	url := fmt.Sprintf("http://%s%s", self.Address(), path)

	req, err := http.NewRequest(method, url, data.(io.Reader))

	if err != nil {
		return err
	}

	mimeType := reqCodec.(interface {
		MimeType() string
	}).MimeType()

	req.Header.Set("Content-Type", mimeType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = getCodec(param.ACKCodecName).Decode(resp.Body, param.ACKMsg)

	if log.IsDebugEnabled() {
		log.Debugf("#http.recv(%s) '%s' %s | [%d] Message(%s) %s",
			self.Name(),
			resp.Request.Method,
			path,
			resp.StatusCode,
			getTypeName(param.ACKMsg),
			cellnet.MessageToString(param.ACKMsg))
	}

	return err
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
