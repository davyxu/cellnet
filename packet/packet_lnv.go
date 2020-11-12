package packet

import (
	"encoding/binary"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/x/io"
	"io"
)

const (
	lnv_bodyHeaderSize = 2  // 包体大小字段
	lnv_nameLen        = 2  // 名字长度
	maxNameLen         = 64 // 最大消息名长度
)

// 接收基于ID格式的封包流程
func RecvLenNameValue(reader io.Reader, maxPacketSize int) (msg interface{}, err error) {

	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, lnv_bodyHeaderSize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(reader, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < lnv_bodyHeaderSize {
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

	msgNameLen := binary.LittleEndian.Uint16(body)

	if msgNameLen > maxNameLen {
		return nil, ErrMaxNameLen
	}

	msgNameBytes := body[lnv_nameLen : lnv_nameLen+msgNameLen]

	msgData := body[lnv_nameLen+msgNameLen:]

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessageByName(string(msgNameBytes), msgData)
	if err != nil {
		return nil, err
	}

	return
}

// 发送基于ID格式的封包流程
func SendLenNameValue(writer io.Writer, ctx cellnet.ContextSet, data interface{}) error {

	var (
		msgData []byte
		msgName string
		meta    *cellnet.MessageMeta
	)

	switch m := data.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = m.MsgData
		msgName = m.MsgName
	default: // 发普通编码包
		var err error

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(data, ctx)

		if err != nil {
			return err
		}

		msgName = meta.FullName()
	}

	bodySize := lnv_nameLen + len(msgName) + len(msgData)
	pkt := make([]byte, lnv_bodyHeaderSize+bodySize)

	pos := pkt[:]

	// Length
	binary.LittleEndian.PutUint16(pos, uint16(bodySize))
	pos = pos[lidv_bodyHeaderSize:]

	// Name
	binary.LittleEndian.PutUint16(pos, uint16(len(msgName)))
	pos = pos[lnv_nameLen:]

	// 拷贝名字
	copy(pos, msgName)
	pos = pos[len(msgName):]

	// Value
	copy(pos, msgData)

	// 将数据写入Socket
	err := xio.WriteFull(writer, pkt)

	// Codec中使用内存池时的释放位置
	if meta != nil {
		codec.FreeCodecResource(meta.Codec, msgData, ctx)
	}

	return err
}
