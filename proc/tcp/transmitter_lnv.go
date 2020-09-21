package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/packet"
	"io"
	"net"
)

type TCPLNVTransmitter struct {
}

func (TCPLNVTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	opt := ses.Peer().(socketOpt)

	if conn, ok := reader.(net.Conn); ok {

		// 有读超时时，设置超时
		opt.ApplySocketReadTimeout(conn, func() {

			msg, err = packet.RecvLenNameValue(reader, opt.MaxPacketSize())

		})
	}

	return
}

func (TCPLNVTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) (err error) {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	opt := ses.Peer().(socketOpt)

	// 有写超时时，设置超时
	opt.ApplySocketWriteTimeout(writer.(net.Conn), func() {

		err = packet.SendLenNameValue(writer, ses.(cellnet.ContextSet), msg)

	})

	return
}
