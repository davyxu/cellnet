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
	errNotHandled = errors.New("Request not handled")
	errNotFound   = errors.New("404 Not found")
)

func (self *httpAcceptor) Start() cellnet.Peer {

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	go http.ListenAndServe(self.Address(), self)

	return self
}

func (self *httpAcceptor) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	ses := newHttpSession(self, req, res)

	var msg interface{}
	var err error

	// 请求转消息，文件处理
	meta := cellnet.HttpMetaByMethodURL(req.Method, req.URL.Path)
	if meta != nil {

		// 直接打开页面时，无需创建消息
		if meta.RequestType != nil {
			msg = reflect.New(meta.RequestType).Interface()

			if err := meta.RequestCodec.Decode(req, msg); err != nil {
				return
			}
		}

	}

	if err == errNotHandled {
		msg, err = self.ServeFileWithDir(res, req)
	}

	if err != nil {

		log.Warnf("#recv %s(%s) %s | 404 NotFound",
			req.Method,
			self.Name(),
			req.URL.Path)

		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(err.Error()))

		return
	}

	// 处理消息及页面下发
	self.PostEvent(&cellnet.RecvMsgEvent{ses, msg})
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
