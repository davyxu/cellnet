package socket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/davyxu/cellnet"
	"io"
	"net"
	"sync"
)

const (
	PackageHeaderSize = 8
	MaxPacketSize     = 1024 * 8
)

// 封包流
type PacketStream interface {
	Read() (*cellnet.Packet, error)
	Write(pkt *cellnet.Packet) error
	Close() error
	Raw() net.Conn
}

type ltvStream struct {
	recvser      uint16
	sendser      uint16
	conn         net.Conn
	sendtagGuard sync.RWMutex
}

var (
	packageTagNotMatch     = errors.New("ReadPacket: package tag not match")
	packageDataSizeInvalid = errors.New("ReadPacket: package crack, invalid size")
	packageTooBig          = errors.New("ReadPacket: package too big")
)

// 参考hub_client.go
// Read a packet from a datastream interface , return packet struct
func (self *ltvStream) Read() (p *cellnet.Packet, err error) {

	headdata := make([]byte, PackageHeaderSize)

	if _, err = io.ReadFull(self.conn, headdata); err != nil {
		return nil, err
	}

	p = &cellnet.Packet{}

	headbuf := bytes.NewReader(headdata)

	// 读取ID
	if err = binary.Read(headbuf, binary.LittleEndian, &p.MsgID); err != nil {
		return nil, err
	}

	// 读取序号
	var ser uint16
	if err = binary.Read(headbuf, binary.LittleEndian, &ser); err != nil {
		return nil, err
	}

	// 读取整包大小
	var fullsize uint16
	if err = binary.Read(headbuf, binary.LittleEndian, &fullsize); err != nil {
		return nil, err
	}

	// 封包太大
	if fullsize > MaxPacketSize {
		return nil, packageTooBig
	}

	// 序列号不匹配
	if self.recvser != ser {
		return nil, packageTagNotMatch
	}

	dataSize := fullsize - PackageHeaderSize
	if dataSize < 0 {
		return nil, packageDataSizeInvalid
	}

	// 读取数据
	p.Data = make([]byte, dataSize)
	if _, err = io.ReadFull(self.conn, p.Data); err != nil {
		return nil, err
	}

	// 增加序列号值
	self.recvser++

	return
}

// Write a packet to datastream interface
func (self *ltvStream) Write(pkt *cellnet.Packet) (err error) {

	outbuff := bytes.NewBuffer([]byte{})

	// 防止将Send放在go内造成的多线程冲突问题
	self.sendtagGuard.Lock()
	defer self.sendtagGuard.Unlock()

	// 写ID
	if err = binary.Write(outbuff, binary.LittleEndian, pkt.MsgID); err != nil {
		return
	}

	// 写序号
	if err = binary.Write(outbuff, binary.LittleEndian, self.sendser); err != nil {
		return
	}

	// 写包大小
	if err = binary.Write(outbuff, binary.LittleEndian, uint16(len(pkt.Data)+PackageHeaderSize)); err != nil {
		return
	}

	// 发包头
	if _, err = self.conn.Write(outbuff.Bytes()); err != nil {
		return err
	}

	// 发包内容
	if _, err = self.conn.Write(pkt.Data); err != nil {
		return err
	}

	// 增加序号值
	self.sendser++

	return
}

func (self *ltvStream) Close() error {
	return self.conn.Close()
}

func (self *ltvStream) Raw() net.Conn {
	return self.conn
}

func NewPacketStream(conn net.Conn) PacketStream {
	return &ltvStream{
		conn:    conn,
		recvser: 1,
		sendser: 1,
	}
}
