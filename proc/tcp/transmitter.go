package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"io"
)

type TCPMessageTransmitter struct {
}

func (TCPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	msg, err = util.RecvLTVPacket(reader)

	return
}

func (TCPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	return util.SendLTVPacket(writer, msg)
}
