package util

import (
	"bytes"
	"encoding/binary"
)

// 封包写入
type BinaryWriter struct {
	buffer bytes.Buffer
}

// 写入的数据字节数
func (p *BinaryWriter) Len() int {
	return p.buffer.Len()
}

// 写入任意值
func (p *BinaryWriter) WriteValue(v interface{}) error {
	return binary.Write(&p.buffer, binary.LittleEndian, v)
}

// 写入的字节数组
func (p *BinaryWriter) Raw() []byte {
	return p.buffer.Bytes()
}

func (p *BinaryWriter) WriteString(v string) error {

	if err := p.WriteValue(uint16(len(v))); err != nil {
		return err
	}

	if err := p.WriteValue([]byte(v)); err != nil {
		return err
	}

	return nil

}
