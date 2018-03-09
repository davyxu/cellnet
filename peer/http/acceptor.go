package http

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"net/http"
	"reflect"
)

type httpAcceptor struct {
	peer.CorePeerProperty
	peer.CoreProcessorBundle
}

var (
	errNotFound = errors.New("404 Not found")
)

func (self *httpAcceptor) Start() cellnet.Peer {

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	go func() {
		err := http.ListenAndServe(self.Address(), self)
		if err != nil {
			log.Errorf("#listen failed(%s) %v", self.NameOrAddress(), err.Error())
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

		log.Warnf("#recv http.%s '%s' %s | [%d] File not found",
			self.Name(),
			req.Method,
			req.URL.Path,
			http.StatusNotFound)

		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(err.Error()))

		return
	}

	if fileHandled {
		log.Debugf("#recv(%s) http.%s %s | [%d] File",
			self.Name(),
			req.Method,
			req.URL.Path,
			http.StatusOK)
		return
	}

	log.Warnf("#recv(%s) http.%s %s | Unhandled",
		self.Name(),
		req.Method,
		req.URL.Path)

	return
OnError:
	log.Errorf("#recv(%s) http.%s %s | [%d] %s",
		self.Name(),
		req.Method,
		req.URL.Path,
		http.StatusInternalServerError,
		err.Error())

	http.Error(ses.resp, err.Error(), http.StatusInternalServerError)
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
