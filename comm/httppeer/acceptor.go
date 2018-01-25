package httppeer

import (
	"github.com/davyxu/cellnet"
	"net/http"
)

type httpAcceptor struct {
	cellnet.CoreTagger
	cellnet.CorePeerInfo
	cellnet.CoreDuplexEventProc
}

func (self *httpAcceptor) Start() cellnet.Peer {

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	go http.ListenAndServe(self.Address(), self)

	return self
}

func (self *httpAcceptor) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	meta := cellnet.MessageMetaByURL(req.URL.Path)
	if meta == nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	msg := meta.NewType()

	if err := meta.Codec.Decode(req, msg); err != nil {
		return
	}

	ses := newHttpSession(self, res)

	self.CallInboundProc(&cellnet.RecvMsgEvent{ses, msg})
}

// 停止侦听器
func (self *httpAcceptor) Stop() {

}

func (self *httpAcceptor) TypeName() string {
	return "http.Acceptor"
}

func init() {

	cellnet.RegisterPeerCreator(func() cellnet.Peer {
		p := &httpAcceptor{}

		return p
	})
}
