package tcptransmit

import (
	"encoding/binary"
	cellevent "github.com/davyxu/cellnet/event"
	"github.com/davyxu/cellnet/peer/tcp"
	celltransmit "github.com/davyxu/cellnet/transmit"
	"io"
)

func RecvMessage(ses *tcp.Session) (ev *cellevent.RecvMsg, err error) {

	if TestEnableRecvPanic {
		panic("emulate recv crash")
	}

	opt := ses.Peer.SocketOption
	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, packetHeaderSize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(ses.Raw(), sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < packetHeaderSize {
		return nil, celltransmit.ErrMinPacket
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint16(sizeBuffer)

	if opt.MaxPacketSize > 0 && int(size) >= opt.MaxPacketSize {
		return nil, celltransmit.ErrMaxPacket
	}

	// 分配包体大小
	body := make([]byte, size)

	// 读取包体数据
	_, err = io.ReadFull(ses.Raw(), body)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(body) < msgIDLen {
		return nil, celltransmit.ErrShortMsgID
	}

	msgid := binary.LittleEndian.Uint16(body)

	msgData := body[msgIDLen:]

	ev = &cellevent.RecvMsg{
		Ses:     ses,
		MsgID:   int(msgid),
		MsgData: msgData,
	}

	return
}
