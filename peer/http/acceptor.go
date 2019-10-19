package http

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	"html/template"
	"net"
	"net/http"
	"strings"
	"time"
)

type httpAcceptor struct {
	peer.CorePeerProperty
	peer.CoreProcBundle
	peer.CoreContextSet

	sv *http.Server

	httpDir  string
	httpRoot string

	templateDir   string
	delimsLeft    string
	delimsRight   string
	templateExts  []string
	templateFuncs []template.FuncMap

	listener net.Listener
}

var (
	errNotFound = errors.New("404 Not found")
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func (self *httpAcceptor) Port() int {
	if self.listener == nil {
		return 0
	}

	return self.listener.Addr().(*net.TCPAddr).Port
}

func (self *httpAcceptor) IsReady() bool {
	return self.Port() != 0
}

func (self *httpAcceptor) WANAddress() string {

	pos := strings.Index(self.Address(), ":")
	if pos == -1 {
		return self.Address()
	}

	host := self.Address()[:pos]

	if host == "" {
		host = util.GetLocalIP()
	}

	return util.JoinAddress(host, self.Port())
}

func (self *httpAcceptor) Start() cellnet.Peer {

	self.sv = &http.Server{Addr: self.Address(), Handler: self}

	ln, err := util.DetectPort(self.Address(), func(a *util.Address, port int) (interface{}, error) {
		return net.Listen("tcp", a.HostPortString(port))
	})

	if err != nil {

		log.Errorf("#http.listen failed(%s) %v", self.Name(), err.Error())

		return self
	}

	self.listener = ln.(net.Listener)

	log.Infof("#http.listen(%s) http://%s", self.Name(), self.WANAddress())

	go func() {

		err = self.sv.Serve(tcpKeepAliveListener{self.listener.(*net.TCPListener)})
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("#http.listen failed(%s) %v", self.Name(), err.Error())
		}

	}()

	return self
}

func (self *httpAcceptor) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	ses := newHttpSession(self, req, res)

	var msg interface{}
	var err error
	var fileHandled bool

	// 处理消息及页面下发
	self.ProcEvent(&cellnet.RecvMsgEvent{Ses: ses, Msg: msg})

	if ses.err != nil {
		err = ses.err
		goto OnError
	}

	if ses.respond {
		return
	}

	// 处理静态文件
	_, err, fileHandled = self.ServeFileWithDir(res, req)

	if err != nil {

		// 或者是普通消息没有Handled
		log.Warnf("#http.recv(%s) '%s' %s | [%d] Not found",
			self.Name(),
			req.Method,
			req.URL.Path,
			http.StatusNotFound)

		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(err.Error()))

		return
	}

	if fileHandled {
		log.Debugf("#http.recv(%s) '%s' %s | [%d] File",
			self.Name(),
			req.Method,
			req.URL.Path,
			http.StatusOK)
		return
	}

	log.Warnf("#http.recv(%s) '%s' %s | Unhandled",
		self.Name(),
		req.Method,
		req.URL.Path)

	return
OnError:
	log.Errorf("#http.recv(%s) '%s' %s | [%d] %s",
		self.Name(),
		req.Method,
		req.URL.Path,
		http.StatusInternalServerError,
		err.Error())

	http.Error(ses.resp, err.Error(), http.StatusInternalServerError)
}

// 停止侦听器
func (self *httpAcceptor) Stop() {

	if err := self.sv.Shutdown(nil); err != nil {
		log.Errorf("#http.stop failed(%s) %v", self.Name(), err.Error())
	}
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
