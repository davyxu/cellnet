package socket

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/davyxu/cellnet"
)

const (
	PackageHeaderSize = 8 // MsgID(uint32) + Ser(uint16) + Size(uint16)
	MaxPacketSize     = 1024 * 8
)

// 封包流
type PacketStream interface {
	Read() (*cellnet.Packet, error)
	Write(pkt *cellnet.Packet) error
	Flush() error
	Close() error
	Raw() net.Conn
}

type ltvStream struct {
	recvser      uint16
	sendser      uint16
	conn         net.Conn
	sendtagGuard sync.RWMutex

	outputWriter     *bufio.Writer
	outputHeadBuffer *bytes.Buffer

	inputHeadBuffer []byte
	headReader      *bytes.Reader
}

var (
	packageTagNotMatch     = errors.New("ReadPacket: package tag not match")
	packageDataSizeInvalid = errors.New("ReadPacket: package crack, invalid size")
	packageTooBig          = errors.New("ReadPacket: package too big")
)

// 从socket读取1个封包,并返回
func (self *ltvStream) Read() (p *cellnet.Packet, err error) {

	if _, err = self.headReader.Seek(0, 0); err != nil {
		return nil, err
	}

	if _, err = io.ReadFull(self.conn, self.inputHeadBuffer); err != nil {
		return nil, err
	}

	p = &cellnet.Packet{}

	// 读取ID
	if err = binary.Read(self.headReader, binary.LittleEndian, &p.MsgID); err != nil {
		return nil, err
	}

	// 读取序号
	var ser uint16
	if err = binary.Read(self.headReader, binary.LittleEndian, &ser); err != nil {
		return nil, err
	}

	// 读取整包大小
	var fullsize uint16
	if err = binary.Read(self.headReader, binary.LittleEndian, &fullsize); err != nil {
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

// 将一个封包发送到socket
func (self *ltvStream) Write(pkt *cellnet.Packet) (err error) {

	// 防止将Send放在go内造成的多线程冲突问题
	self.sendtagGuard.Lock()
	defer self.sendtagGuard.Unlock()

	self.outputHeadBuffer.Reset()

	// 写ID
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, pkt.MsgID); err != nil {
		return err
	}

	// 写序号
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, self.sendser); err != nil {
		return err
	}

	// 写包大小
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, uint16(len(pkt.Data)+PackageHeaderSize)); err != nil {
		return err
	}

	// 发包头
	if _, err = self.outputWriter.Write(self.outputHeadBuffer.Bytes()); err != nil {
		return err
	}

	// 发包内容
	if _, err = self.outputWriter.Write(pkt.Data); err != nil {
		return err
	}

	// 增加序号值
	self.sendser++

	return
}

func (self *ltvStream) Flush() error {

	if err := self.outputWriter.Flush(); err != nil && err != io.ErrShortWrite {
		return err
	}

	return nil
}

// 关闭
func (self *ltvStream) Close() error {
	return self.conn.Close()
}

// 裸socket操作
func (self *ltvStream) Raw() net.Conn {
	return self.conn
}

// 封包流 relay模式: 在封包头有clientid信息
func NewPacketStream(conn net.Conn) PacketStream {

	s := &ltvStream{
		conn:             conn,
		recvser:          1,
		sendser:          1,
		outputWriter:     bufio.NewWriter(conn),
		outputHeadBuffer: bytes.NewBuffer([]byte{}),
		inputHeadBuffer:  make([]byte, PackageHeaderSize),
	}
	s.headReader = bytes.NewReader(s.inputHeadBuffer)

	return s
}
