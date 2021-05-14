package tcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	cellcodec "github.com/davyxu/cellnet/codec"
	cellevent "github.com/davyxu/cellnet/event"
	xbytes "github.com/davyxu/x/bytes"
	xio "github.com/davyxu/x/io"
	"io"
)

var (
	ErrMaxPacket  = errors.New("packet over size")
	ErrMinPacket  = errors.New("packet short size")
	ErrShortMsgID = errors.New("short msgid")
)

const (
	packetHeaderSize = 2 // 包体大小字段
	msgIDLen         = 2 // 消息ID字段
)

type NoCryptMessage struct {
	Msg interface{}
}

func (self *NoCryptMessage) Message() interface{} {
	return self.Msg
}

var (
	TestEnableRecvPanic bool
	TestEnableSendPanic bool
)

func RecvMessage(ses *Session) (ev *cellevent.RecvMsg, err error) {

	if TestEnableRecvPanic {
		panic("emulate recv crash")
	}

	opt := ses.Peer.SocketOption
	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, packetHeaderSize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(ses.conn, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < packetHeaderSize {
		return nil, ErrMinPacket
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint16(sizeBuffer)

	if opt.MaxPacketSize > 0 && int(size) >= opt.MaxPacketSize {
		return nil, ErrMaxPacket
	}

	// 分配包体大小
	body := make([]byte, size)

	// 读取包体数据
	_, err = io.ReadFull(ses.conn, body)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(body) < msgIDLen {
		return nil, ErrShortMsgID
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

func SendMessage(ses *Session, ev *cellevent.SendMsg) error {

	if TestEnableSendPanic {
		panic("emulate send crash")
	}

	ps := &ses.Peer.PropertySet

	var (
		msgData []byte
	)

	if ev.MessageData() != nil {
		msgData = ev.MessageData()
	} else if raw, ok := ev.Message().(*NoCryptMessage); ok {
		data, meta, err := cellcodec.Encode(raw.Msg, ps)
		if err != nil {
			return fmt.Errorf("encode msg failed, %+v", raw.Msg)
		} else {
			msgData = data
			ev.MsgID = meta.ID
		}
	} else {
		panic(fmt.Sprintf("invalid message %+v", ev.Message()))
	}

	bodySize := msgIDLen + len(msgData)
	composeBuffer := make([]byte, packetHeaderSize+bodySize)
	writer := xbytes.NewWriter(composeBuffer)

	writer.WriteUint16(uint16(bodySize))
	writer.WriteUint16(uint16(ev.MsgID))
	writer.Write(msgData)

	// 将数据写入Socket

	return xio.WriteFull(ses.conn, composeBuffer)
}
