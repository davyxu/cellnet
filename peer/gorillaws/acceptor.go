package gorillaws

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

type wsAcceptor struct {
	peer.CoreSessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreProcBundle

	certfile string
	keyfile  string

	upgrader websocket.Upgrader
	// 保存端口
	listener net.Listener

	sv *http.Server
}

func (self *wsAcceptor) SetUpgrader(upgrader interface{}) {
	self.upgrader = upgrader.(websocket.Upgrader)
}
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

	var (
		addrObj *util.Address
		err     error
		raw     interface{}
	)

	raw, err = util.DetectPort(self.Address(), func(a *util.Address, port int) (interface{}, error) {
		addrObj = a
		return net.Listen("tcp", a.HostPortString(port))
	})

	if err != nil {
		log.Errorf("#ws.listen failed(%s) %v", self.Name(), err.Error())
		return self
	}

	self.listener = raw.(net.Listener)

	mux := http.NewServeMux()

	if addrObj.Path == "" {
		addrObj.Path = "/"
	}

	mux.HandleFunc(addrObj.Path, func(w http.ResponseWriter, r *http.Request) {

		c, err := self.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Debugln(err)
			return
		}

		ses := newSession(c, self, nil)
		ses.SetContext("request", r)
		ses.Start()

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: ses, Msg: &cellnet.SessionAccepted{}})

	})

	self.sv = &http.Server{Addr: addrObj.HostPortString(self.Port()), Handler: mux}

	go func() {

		log.Infof("#ws.listen(%s) %s", self.Name(), addrObj.String(self.Port()))

		if self.certfile != "" && self.keyfile != "" {
			err = self.sv.ServeTLS(self.listener, self.certfile, self.keyfile)
		} else {
			err = self.sv.Serve(self.listener)
		}

		if err != nil {
			log.Errorf("#ws.listen. failed(%s) %v", self.Name(), err.Error())
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
		p := &wsAcceptor{
			upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		}

		return p
	})
}
