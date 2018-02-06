package kcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/util"
	"io"
)

func (self *kcpContext) input(data []byte) {

	log.Debugln("input", self.ses.Peer().(cellnet.PeerProperty).Name(), len(data), data)

	if ret := self.kcp.Input(data, true, true); ret != 0 {
		log.Errorln("kcp input ret: ", ret)
	}

	self.readSignal <- struct{}{}
}

func (self *kcpContext) Read(b []byte) (n int, err error) {

	defer log.Debugln("read", self.ses.Peer().(cellnet.PeerProperty).Name(), n, b)

	for {

		if len(self.bufptr) > 0 { // copy from buffer into b
			n = copy(b, self.bufptr)
			self.bufptr = self.bufptr[n:]
			return n, nil
		}

		if size := self.kcp.PeekSize(); size > 0 {

			if len(b) >= size { // direct write to b
				self.kcp.Recv(b)
				return size, nil

			}

			// resize kcp receive buffer
			// to make sure recvbuf has enough capacity
			if cap(self.recvbuf) < size {
				self.recvbuf = make([]byte, size)
			}

			// resize recvbuf slice length
			self.recvbuf = self.recvbuf[:size]
			self.kcp.Recv(self.recvbuf)
			n = copy(b, self.recvbuf)      // copy to b
			self.bufptr = self.recvbuf[n:] // update pointer
			return n, nil
		}

		if self.closed {
			return 0, io.EOF
		}

		<-self.readSignal
	}
}

func (self *kcpContext) RecvLTVPacket() (msg interface{}, err error) {

	pktReader, err := util.RecvVariableLengthPacket(self)
	if err != nil {
		log.Errorln(err)
		return
	}

	// 读取消息ID
	var msgid uint16
	if err = pktReader.ReadValue(&msgid); err != nil {
		return
	}

	msgData := pktReader.RemainBytes()

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(int(msgid), msgData)

	return
}
