package http

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"html/template"
	"net/http"
)

type RequestMatcher interface {
	Match(method, url string) bool
}
type RespondProc interface {
	WriteRespond(*httpSession) error
}

var (
	ErrUnknownOperation = errors.New("Unknown http operation")
)

type httpSession struct {
	peer.CoreContextSet
	*peer.CoreProcBundle
	req  *http.Request
	resp http.ResponseWriter

	// 单独保存的保存Peer接口
	peerInterface cellnet.Peer

	t *template.Template

	respond bool
	err     error
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

// 发送封包
func (self *httpSession) Send(raw interface{}) {

	if proc, ok := raw.(RespondProc); ok {
		self.err = proc.WriteRespond(self)
		self.respond = true
	} else {
		self.err = ErrUnknownOperation
	}

}

func newHttpSession(acc *httpAcceptor, req *http.Request, response http.ResponseWriter) *httpSession {

	return &httpSession{
		req:            req,
		resp:           response,
		peerInterface:  acc,
		t:              acc.Compile(),
		CoreProcBundle: acc.GetBundle(),
	}
}
