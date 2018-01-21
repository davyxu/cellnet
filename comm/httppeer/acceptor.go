package httppeer

import (
	"github.com/davyxu/cellnet"
	"net/http"
	"strings"
)

type httpAcceptor struct {
	cellnet.CoreTagger
	cellnet.CorePeerInfo
}

func (self *httpAcceptor) Start() cellnet.Peer {

	go http.ListenAndServe(self.Address(), self)

	return self
}

func (self *httpAcceptor) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")

	if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" || contentType != "" {
		if strings.Contains(contentType, "form-urlencoded") {
			//context.Invoke(Form(obj, ifacePtr...))
		} else if strings.Contains(contentType, "multipart/form-data") {
			//context.Invoke(MultipartForm(obj, ifacePtr...))
		} else if strings.Contains(contentType, "json") {
			//context.Invoke(Json(obj, ifacePtr...))
		} else {
			//var errors Errors
			//if contentType == "" {
			//	errors.Add([]string{}, ContentTypeError, "Empty Content-Type")
			//} else {
			//	errors.Add([]string{}, ContentTypeError, "Unsupported Content-Type")
			//}
			//context.Map(errors)
		}
	} else {
		//context.Invoke(Form(obj, ifacePtr...))
	}

	meta := cellnet.MessageMetaByName(req.URL.Path)
	if meta == nil {
		return
	}
}

// 停止侦听器
func (self *httpAcceptor) Stop() {

}

func (self *httpAcceptor) TypeName() string {
	return "http.Acceptor"
}

func init() {

	cellnet.RegisterPeerCreator(func() cellnet.Peer {
		p := &httpAcceptor{}

		return p
	})
}
