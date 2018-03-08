package http

import "github.com/davyxu/cellnet"

type StatusCode int

func (self StatusCode) WriteRespond(ses *httpSession) error {

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	log.Debugf("#recv(%s) http.%s %s | [%d] Status",
		peerInfo.Name(),
		ses.req.Method,
		ses.req.URL.Path,
		self)

	ses.resp.WriteHeader(int(self))
	return nil
}
