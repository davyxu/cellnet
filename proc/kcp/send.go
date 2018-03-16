package kcp

import (
	"github.com/davyxu/cellnet/peer/udp"
)

func (self *kcpContext) output(data []byte) {

	//log.Debugln("output", self.ses.Peer().(cellnet.PeerProperty).Name(), len(data), data)

	writer := self.ses.(udp.DataWriter)

	writer.WriteData(data)
}

func (self *kcpContext) Write(p []byte) (n int, err error) {

	//log.Debugln("write", self.ses.Peer().(cellnet.PeerProperty).Name(), len(p), p)

	self.kcp.Send(p)

	return len(p), nil
}
