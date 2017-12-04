package packet

import (
	"bytes"
	"encoding/binary"
)

type PacketReader struct {
	raw []byte

	reader *bytes.Reader
}

// 初始化缓冲，清空读取器
func (p *PacketReader) Init(raw []byte) {
	p.raw = raw
	p.reader = nil
}

func (p *PacketReader) Raw() []byte {
	return p.raw
}

func (p *PacketReader) prepareReader() {

	if p.reader == nil {
		p.reader = bytes.NewReader(p.raw)
	}

}

// 未读取的数据字节数
func (p *PacketReader) RemainLen() int {

	p.prepareReader()

	return p.reader.Len()
}

// 剩下的未读取的字节
func (p *PacketReader) RemainBytes() []byte {

	p.prepareReader()

	return p.raw[len(p.raw)-p.reader.Len():]
}

// 从字节数组中读取值
func (p *PacketReader) ReadValue(v interface{}) error {

	p.prepareReader()

	return binary.Read(p.reader, binary.LittleEndian, v)
}

// 读取字符串
func (p *PacketReader) ReadString(str *string) error {
	// 读取字符串长度
	var size uint16
	if err := p.ReadValue(&size); err != nil {
		return err
	}

	// 分配字符串空间
	body := make([]byte, size)

	// 读取字符串值
	if err := p.ReadValue(&body); err != nil {
		return err
	}

	// 返回字符串
	*str = string(body)

	return nil
}
