package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"html/template"
	"net/http"
)

type RequestMatcher interface {
	Match(method, url string) bool
}

type httpSession struct {
	peer.CorePropertySet
	*peer.CoreProcessorBundle
	req  *http.Request
	resp http.ResponseWriter

	// 单独保存的保存Peer接口
	peerInterface cellnet.Peer

	t *template.Template
}

func (self *httpSession) Match(method, url string) bool {

	return self.req.Method == method && self.req.URL.Path == url
}

func (self *httpSession) Request() *http.Request {
	return self.req
}

func (self *httpSession) Response() http.ResponseWriter {
	return self.resp
}

// 取原始连接
func (self *httpSession) Raw() interface{} {
	return nil
}

func (self *httpSession) ID() int64 {
	return 0
}

// 取原始连接
func (self *httpSession) Close() {
}

// 取会话归属的通讯端
func (self *httpSession) Peer() cellnet.Peer {
	return self.peerInterface
}

type StatusCode int

type HTML struct {
	Code int

	PageTemplate string

	TemplateModel interface{}
}

// 发送封包
func (self *httpSession) Send(raw interface{}) {

	switch msg := raw.(type) {
	case StatusCode:
		self.resp.WriteHeader(int(msg))
	case HTML:
		writeHTMLRespond(self, msg.Code, msg.PageTemplate, msg.TemplateModel)
	default:
		writeMessageRespond(self, msg)
	}
}

func newHttpSession(peerIns cellnet.Peer, req *http.Request, response http.ResponseWriter) cellnet.Session {

	return &httpSession{
		req:           req,
		resp:          response,
		peerInterface: peerIns,
		t:             compile(prepareOptions([]Options{})),
		CoreProcessorBundle: peerIns.(interface {
			GetBundle() *peer.CoreProcessorBundle
		}).GetBundle(),
	}
}
