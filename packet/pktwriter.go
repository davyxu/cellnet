package packet

import (
	"bytes"
	"encoding/binary"
)

// 封包写入
type PacketWriter struct {
	buffer bytes.Buffer
}

// 写入的数据字节数
func (p *PacketWriter) Len() uint16 {
	return uint16(p.buffer.Len())
}

// 写入任意值
func (p *PacketWriter) WriteValue(v interface{}) error {
	return binary.Write(&p.buffer, binary.LittleEndian, v)
}

// 写入的字节数组
func (p *PacketWriter) Raw() []byte {
	return p.buffer.Bytes()
}

func (p *PacketWriter) WriteString(v string) error {

	if err := p.WriteValue(uint16(len(v))); err != nil {
		return err
	}

	if err := p.WriteValue([]byte(v)); err != nil {
		return err
	}

	return nil

}
