package httpfile

import (
	"github.com/davyxu/cellnet"
	"net/http"
)

// 接受器
type httpFileAcceptor struct {
	cellnet.CorePeerInfo
	cellnet.CoreTagger
}

func (self *httpFileAcceptor) Start() cellnet.Peer {

	// 注册消息
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	go http.ListenAndServe(self.Address(), nil)

	return self
}

func (self *httpFileAcceptor) Stop() {

	// TODO graceful shutdown
}

func (self *httpFileAcceptor) TypeName() string {
	return "http.file.Acceptor"
}

func init() {

	cellnet.RegisterPeerCreator(func() cellnet.Peer {
		p := &httpFileAcceptor{}

		return p
	})
}
