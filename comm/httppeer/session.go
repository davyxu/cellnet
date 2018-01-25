package httppeer

import (
	"github.com/davyxu/cellnet"
	"net/http"
)

type Status struct {
	Code int
}

type httpSession struct {
	cellnet.CoreTagger
	Response http.ResponseWriter

	// 单独保存的保存Peer接口
	peerInterface cellnet.Peer
}

// 取原始连接
func (self *httpSession) Raw() interface{} {
	return self.Response
}

func (self *httpSession) ID() int64 {
	return 0
}

// 取会话归属的通讯端
func (self *httpSession) Peer() cellnet.Peer {
	return self.peerInterface
}

func (self *httpSession) Close() {
}

// 发送封包
func (self *httpSession) Send(raw interface{}) {

	switch msg := raw.(type) {
	case Status:
		self.Response.WriteHeader(msg.Code)
	default:

		data, _, err := cellnet.EncodeMessage(msg)
		if err != nil {
			self.Response.WriteHeader(http.StatusNotFound)
			return
		}

		self.Response.Write(data)
	}
}

func newHttpSession(peer cellnet.Peer, response http.ResponseWriter) cellnet.Session {

	return &httpSession{
		peerInterface: peer,
		Response:      response,
	}
}
