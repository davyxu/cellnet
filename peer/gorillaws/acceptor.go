package gorillaws

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/ulog"
	xframe "github.com/davyxu/x/frame"
	"github.com/davyxu/x/net"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

type wsAcceptor struct {
	peer.CoreSessionManager
	peer.CorePeerProperty
	xframe.PropertySet
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
		addrObj *xnet.Address
		err     error
		raw     interface{}
	)

	raw, err = xnet.DetectPort(self.Address(), func(a *xnet.Address, port int) (interface{}, error) {
		addrObj = a
		return net.Listen("tcp", a.HostPortString(port))
	})

	if err != nil {
		ulog.Errorf("#ws.listen failed(%s) %v", self.Name(), err.Error())
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
			ulog.Debugln(err)
			return
		}

		ses := newSession(c, self, nil)
		ses.Set("request", r)
		ses.Start()

		self.ProcEvent(cellnet.BuildSystemEvent(ses, &cellnet.SessionAccepted{}))
	})

	self.sv = &http.Server{Addr: addrObj.HostPortString(self.Port()), Handler: mux}

	go func() {

		ulog.Infof("#ws.listen(%s) %s", self.Name(), addrObj.String(self.Port()))

		if self.certfile != "" && self.keyfile != "" {
			err = self.sv.ServeTLS(self.listener, self.certfile, self.keyfile)
		} else {
			err = self.sv.Serve(self.listener)
		}

		if err != nil {
			ulog.Errorf("#ws.listen. failed(%s) %v", self.Name(), err.Error())
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
