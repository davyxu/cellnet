package kcp

import (
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/util"
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

func (ctx *kcpContext) sendMessage(msg interface{}) error {

	msglog.WriteSendLogger(log, "kcp", ctx.ses, msg)

	// 将用户数据转换为字节数组和消息ID
	data, meta, err := codec.EncodeMessage(msg)

	if err != nil {
		log.Errorf("send message encode error: %s", err)
		return err
	}

	// 创建封包写入器
	var pktWriter util.BinaryWriter

	// 写入消息ID
	if err := pktWriter.WriteValue(uint16(meta.ID)); err != nil {
		return err
	}

	// 写入序列化好的消息数据
	if err := pktWriter.WriteValue(data); err != nil {
		return err
	}

	util.SendVariableLengthPacket(ctx, pktWriter)

	return nil
}
