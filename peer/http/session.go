package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"net/http"
)

type httpSession struct {
	peer.CorePropertySet
	*peer.CoreProcessorBundle
	req  *http.Request
	resp http.ResponseWriter

	// 单独保存的保存Peer接口
	peerInterface cellnet.Peer
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

// 取会话归属的通讯端
func (self *httpSession) Peer() cellnet.Peer {
	return self.peerInterface
}

// 发送封包
func (self *httpSession) Send(raw interface{}) {

	self.SendMessage(&cellnet.SendMsgEvent{self, raw})
}

func newHttpSession(peerIns cellnet.Peer, req *http.Request, response http.ResponseWriter) cellnet.BaseSession {

	return &httpSession{
		req:           req,
		resp:          response,
		peerInterface: peerIns,
		CoreProcessorBundle: peerIns.(interface {
			GetBundle() *peer.CoreProcessorBundle
		}).GetBundle(),
	}
}
