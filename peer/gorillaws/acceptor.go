package gorillaws

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

type wsAcceptor struct {
	peer.CoreSessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle

	certfile string
	keyfile  string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func (self *wsAcceptor) SetHttps(certfile, keyfile string) {

	self.certfile = certfile
	self.keyfile = keyfile
}

func (self *wsAcceptor) Start() cellnet.Peer {

	if self.IsRunning() {
		return self
	}

	self.SetRunning(true)

	urlObj, err := url.Parse(self.Address())

	if err != nil {
		log.Errorf("#websocket.urlparse failed(%s) %v", self.NameOrAddress(), err.Error())
		return self
	}

	if urlObj.Path == "" {
		log.Errorln("#websocket start failed, expect path in url to listen", self.NameOrAddress())
		return self
	}

	mux := http.NewServeMux()

	mux.HandleFunc(urlObj.Path, func(w http.ResponseWriter, r *http.Request) {

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Debugln(err)
			return
		}

		ses := newSession(c, self, nil)

		ses.Start()

		self.PostEvent(&cellnet.RecvMsgEvent{ses, &cellnet.SessionAccepted{}})

	})

	go func() {

		log.Infof("#websocket.listen(%s) %s", self.Name(), self.Address())

		if urlObj.Scheme == "https" {
			err = http.ListenAndServeTLS(urlObj.Host, self.certfile, self.keyfile, mux)
		} else {
			err = http.ListenAndServe(urlObj.Host, mux)
		}

		if err != nil {
			log.Errorf("#websocket.listen. failed(%s) %v", self.NameOrAddress(), err.Error())
		}

		self.SetRunning(false)

	}()

	return self
}

func (self *wsAcceptor) Stop() {

	// TODO 关闭处理
}

func (self *wsAcceptor) TypeName() string {
	return "gorillaws.Acceptor"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &wsAcceptor{}

		return p
	})
}
