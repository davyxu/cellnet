package http

import "github.com/davyxu/cellnet"

type StatusRespond struct {
	StatusCode int
}

func (self *StatusRespond) WriteRespond(ses *httpSession) error {

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	log.Debugf("#recv(%s) http.%s %s | [%d] Status",
		peerInfo.Name(),
		ses.req.Method,
		ses.req.URL.Path,
		self.StatusCode)

	ses.resp.WriteHeader(int(self.StatusCode))
	return nil
}
