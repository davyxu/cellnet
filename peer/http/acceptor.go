package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"net/http"
)

type httpAcceptor struct {
	peer.CorePeerProperty
	peer.CoreProcessorBundle
}

func (self *httpAcceptor) Start() cellnet.Peer {

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	go http.ListenAndServe(self.Address(), self)

	return self
}

func (self *httpAcceptor) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	ses := newHttpSession(self, req, res)

	msg, err := self.ReadMessage(ses)

	if err != nil {

		log.Warnf("#recv %s(%s) %s | 404 NotFound",
			req.Method,
			self.Name(),
			req.URL.Path)

		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(err.Error()))

		return
	}

	if msg != nil {
		self.PostEvent(&cellnet.RecvMsgEvent{ses, msg})
	}
}

// 停止侦听器
func (self *httpAcceptor) Stop() {

}

func (self *httpAcceptor) TypeName() string {
	return "http.Acceptor"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &httpAcceptor{}

		return p
	})
}
