package peer

import (
	"errors"
	"github.com/davyxu/cellnet"
	"io"
	"net"
)

type BundleSupport interface {
	ProcEvent(ev cellnet.Event)

	ReadMessage(ses cellnet.Session) (ev cellnet.Event, err error)

	SendMessage(ev cellnet.Event)
}

type CoreProcBundle struct {
	transmit cellnet.MessageTransmitter
	hooker   cellnet.EventHooker
	callback cellnet.EventCallback
}

func (self *CoreProcBundle) GetBundle() *CoreProcBundle {
	return self
}

func (self *CoreProcBundle) SetTransmitter(v cellnet.MessageTransmitter) {
	self.transmit = v
}

func (self *CoreProcBundle) SetHooker(v cellnet.EventHooker) {
	self.hooker = v
}

func (self *CoreProcBundle) SetCallback(v cellnet.EventCallback) {
	self.callback = v
}

var notHandled = errors.New("Processor: Transimitter nil")

func (self *CoreProcBundle) ReadMessage(ses cellnet.Session) (ev cellnet.Event, err error) {

	if self.transmit != nil {

		opt := ses.Peer().(TCPSocketOptionApply)

		reader, ok := ses.Raw().(io.Reader)

		// 转换错误，或者连接已经关闭时退出
		if !ok || reader == nil {
			return nil, nil
		}

		if conn, ok := reader.(net.Conn); ok {

			if opt.BeginApplyReadTimeout(conn) {
				ev, err = self.transmit.OnRecvMessage(ses)
				opt.EndApplyTimeout(conn)
			} else {
				ev, err = self.transmit.OnRecvMessage(ses)
			}
		}

		return
	}

	return nil, notHandled
}

func (self *CoreProcBundle) SendMessage(ev cellnet.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnOutboundEvent(ev)
	}

	if self.transmit != nil && ev != nil {

		writer, ok := ev.Session().Raw().(io.Writer)

		// 转换错误，或者连接已经关闭时退出
		if ok && writer != nil {
			opt := ev.Session().Peer().(TCPSocketOptionApply)

			conn := writer.(net.Conn)
			if opt.BeginApplyWriteTimeout(conn) {
				self.transmit.OnSendMessage(ev)
				opt.EndApplyTimeout(conn)
			} else {
				self.transmit.OnSendMessage(ev)
			}
		}

	}
}

func (self *CoreProcBundle) ProcEvent(ev cellnet.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnInboundEvent(ev)
	}

	if self.callback != nil && ev != nil {
		self.callback(ev)
	}
}
