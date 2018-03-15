package http

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"html/template"
	"net/http"
	"reflect"
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
}

var (
	errNotFound = errors.New("404 Not found")
)

func (self *httpAcceptor) Start() cellnet.Peer {

	log.Infof("#http.listen(%s) %s", self.Name(), self.Address())

	self.sv = &http.Server{Addr: self.Address(), Handler: self}

	go func() {

		err := self.sv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("#http.listen failed(%s) %v", self.NameOrAddress(), err.Error())
		}

	}()

	return self
}

func (self *httpAcceptor) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	ses := newHttpSession(self, req, res)

	var msg interface{}
	var err error
	var fileHandled bool

	// 请求转消息，文件处理
	meta := cellnet.HttpMetaByMethodURL(req.Method, req.URL.Path)
	if meta != nil {

		// 直接打开页面时，无需创建消息
		if meta.RequestType != nil {
			msg = reflect.New(meta.RequestType).Interface()

			err = meta.RequestCodec.Decode(req, msg)
		}
	}

	if err != nil {
		goto OnError
	}

	// 处理消息及页面下发
	self.PostEvent(&cellnet.RecvMsgEvent{ses, msg})

	if ses.err != nil {
		err = ses.err
		goto OnError
	}

	if ses.responed {
		return
	}

	// 处理静态文件
	msg, err, fileHandled = self.ServeFileWithDir(res, req)

	if err != nil {

		log.Warnf("#http.recv(%s) '%s' %s | [%d] File not found",
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
		log.Errorf("#http.stop failed(%s) %v", self.NameOrAddress(), err.Error())
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
