package gorillaws

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/url"
)

type wsAcceptor struct {
	peer.CoreSessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreProcBundle

	certfile string
	keyfile  string

	// 保存端口
	listener net.Listener

	sv *http.Server
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func (self *wsAcceptor) Port() int {
	if self.listener == nil {
		return 0
	}

	return self.listener.Addr().(*net.TCPAddr).Port
}

func (self *wsAcceptor) IsReady() bool {

	return self.Port() != 0
}

func (self *wsAcceptor) SetHttps(certfile, keyfile string) {

	self.certfile = certfile
	self.keyfile = keyfile
}

func (self *wsAcceptor) Start() cellnet.Peer {

	var addrURL *url.URL
	var err error
	var raw interface{}
	raw, err = util.DetectPort(self.Address(), func(s string) (interface{}, error) {

		addrURL, err = url.Parse(s)

		if err != nil {
			return nil, err
		}

		if addrURL.Path == "" {
			return nil, errors.New("expect path in url to listen")
		}
		return net.Listen("tcp", addrURL.Host)
	})

	if err != nil {
		log.Errorf("#websocket.Listen failed(%s) %v", self.Name(), err.Error())
		return self
	}

	self.listener = raw.(net.Listener)

	mux := http.NewServeMux()

	mux.HandleFunc(addrURL.Path, func(w http.ResponseWriter, r *http.Request) {

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Debugln(err)
			return
		}

		ses := newSession(c, self, nil)

		ses.Start()

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: ses, Msg: &cellnet.SessionAccepted{}})

	})

	self.sv = &http.Server{Addr: addrURL.Host, Handler: mux}

	go func() {

		log.Infof("#websocket.listen(%s) %s", self.Name(), addrURL.String())

		if self.certfile != "" && self.keyfile != "" {
			err = self.sv.ServeTLS(self.listener, self.certfile, self.keyfile)
		} else {
			err = self.sv.Serve(self.listener)
		}

		if err != nil {
			log.Errorf("#websocket.listen. failed(%s) %v", self.Name(), err.Error())
		}

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
