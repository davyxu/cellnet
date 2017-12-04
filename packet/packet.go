package packet

import (
	"encoding/binary"
	"io"
)

const LengthSize = 2

// 接收变长封包
func RecvVariableLengthPacket(inputStream io.Reader) (pktReader PacketReader, err error) {

	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, LengthSize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(inputStream, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint16(sizeBuffer)

	// 分配包体大小
	body := make([]byte, size)

	// 读取包体数据
	_, err = io.ReadFull(inputStream, body)

	// 初始化封包读取器
	pktReader.Init(body)

	return
}

// 发送变长封包
func SendVariableLengthPacket(outputStream io.Writer, pktWriter PacketWriter) error {

	buffer := make([]byte, pktWriter.Len()+LengthSize)

	// 将包体长度写入缓冲
	binary.LittleEndian.PutUint16(buffer, pktWriter.Len())

	// 将包体数据写入缓冲
	copy(buffer[LengthSize:], pktWriter.Raw())

	// 将数据写入Socket
	if _, err := outputStream.Write(buffer); err != nil {
		return err
	}

	return nil
}
