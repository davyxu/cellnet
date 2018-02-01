package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"net/http"
)

type httpAcceptor struct {
	peer.CoreTagger
	proc.CoreDuplexEventProc
	peer.CommunicateConfig

	*StaticFile
}

func (self *httpAcceptor) Start() cellnet.Peer {

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	go http.ListenAndServe(self.Address(), self)

	return self
}

func (self *httpAcceptor) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	meta := cellnet.MessageMetaByHTTPRequest(req.Method, req.URL.Path)
	if meta != nil {

		msg := meta.NewType()

		if err := meta.Codec.Decode(req, msg); err != nil {
			return
		}

		ses := newHttpSession(self, res)

		self.CallInboundProc(&cellnet.RecvMsgEvent{ses, msg})

	} else {

		self.ServeFile(res, req)
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
		p := &httpAcceptor{
			StaticFile: newStaticFile("", "."),
		}

		return p
	})
}
