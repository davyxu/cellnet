package tcptransport

import (
	"encoding/binary"
	cellevent "github.com/davyxu/cellnet/event"
	"github.com/davyxu/cellnet/peer/tcp"
	celltransport "github.com/davyxu/cellnet/transport"
	"io"
)

func RecvMessage(ses *tcp.Session) (ev *cellevent.RecvMsg, err error) {

	if TestEnableRecvPanic {
		panic("emulate recv crash")
	}

	opt := ses.Peer.SocketOption
	// Size为uint32，占4字节
	var sizeBuffer = make([]byte, packetHeaderSize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(ses.Raw(), sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < packetHeaderSize {
		return nil, celltransport.ErrMinPacket
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint32(sizeBuffer)

	if opt.MaxPacketSize > 0 && int(size) >= opt.MaxPacketSize {
		return nil, celltransport.ErrMaxPacket
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
		return nil, celltransport.ErrShortMsgID
	}

	msgid := binary.LittleEndian.Uint32(body)

	msgData := body[msgIDLen:]

	ev = &cellevent.RecvMsg{
		Ses:     ses,
		MsgID:   int(msgid),
		MsgData: msgData,
	}

	return
}
