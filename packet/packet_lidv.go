package packet

import (
	"encoding/binary"
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	xio "github.com/davyxu/x/io"
	"io"
)

var (
	ErrMaxPacket  = errors.New("packet over size")
	ErrMinPacket  = errors.New("packet short size")
	ErrShortMsgID = errors.New("short msgid")
	ErrMaxNameLen = errors.New("max name len")
)

const (
	lidv_bodyHeaderSize = 2 // 包体大小字段
	lidv_msgIDSize      = 2 // 消息ID字段
)

// 接收基于ID格式的封包流程
func RecvLenIDValue(reader io.Reader, maxPacketSize int) (msg interface{}, err error) {

	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, lidv_bodyHeaderSize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(reader, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < lidv_bodyHeaderSize {
		return nil, ErrMinPacket
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint16(sizeBuffer)

	if maxPacketSize > 0 && int(size) >= maxPacketSize {
		return nil, ErrMaxPacket
	}

	// 分配包体大小
	body := make([]byte, size)

	// 读取包体数据
	_, err = io.ReadFull(reader, body)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(body) < lidv_msgIDSize {
		return nil, ErrShortMsgID
	}

	msgid := binary.LittleEndian.Uint16(body)

	msgData := body[lidv_msgIDSize:]

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(int(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}

// 发送基于ID格式的封包流程
func SendLenIDValue(writer io.Writer, ctx cellnet.ContextSet, data interface{}) error {

	var (
		msgData []byte
		msgID   int
		meta    *cellnet.MessageMeta
	)

	switch m := data.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	default: // 发普通编码包
		var err error

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(data, ctx)

		if err != nil {
			return err
		}

		msgID = meta.ID
	}

	pkt := make([]byte, lidv_bodyHeaderSize+lidv_msgIDSize+len(msgData))

	// Length
	binary.LittleEndian.PutUint16(pkt, uint16(lidv_msgIDSize+len(msgData)))

	// Type
	binary.LittleEndian.PutUint16(pkt[lidv_bodyHeaderSize:], uint16(msgID))

	// Value
	copy(pkt[lidv_bodyHeaderSize+lidv_msgIDSize:], msgData)

	// 将数据写入Socket
	err := xio.WriteFull(writer, pkt)

	// Codec中使用内存池时的释放位置
	if meta != nil {
		codec.FreeCodecResource(meta.Codec, msgData, ctx)
	}

	return err
}
