package socket

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/davyxu/cellnet"
	"io"
	"net"
	"sync"
)

const (
	PackageHeaderSize = 8 // MsgID(uint32) + Ser(uint16) + Size(uint16)
)

type TLVStream struct {
	recvser      uint16
	sendser      uint16
	conn         net.Conn
	sendtagGuard sync.RWMutex

	outputWriter     *bufio.Writer
	outputHeadBuffer *bytes.Buffer

	inputHeadBuffer []byte
	headReader      *bytes.Reader

	maxPacketSize int
}

var (
	ErrPackageTagNotMatch     = errors.New("ReadPacket: package tag not match")
	ErrPackageDataSizeInvalid = errors.New("ReadPacket: package crack, invalid size")
	ErrPackageTooBig          = errors.New("ReadPacket: package too big")
)

func (self *TLVStream) SetMaxPacketSize(size int) {
	self.maxPacketSize = size
}

// 从socket读取1个封包,并返回
func (self *TLVStream) Read() (msgid uint32, data []byte, err error) {

	if _, err = self.headReader.Seek(0, 0); err != nil {
		return
	}

	if _, err = io.ReadFull(self.conn, self.inputHeadBuffer); err != nil {
		return
	}

	// 读取ID
	if err = binary.Read(self.headReader, binary.LittleEndian, &msgid); err != nil {
		return
	}

	// 读取序号
	var ser uint16
	if err = binary.Read(self.headReader, binary.LittleEndian, &ser); err != nil {
		return
	}

	// 读取整包大小
	var fullsize uint16
	if err = binary.Read(self.headReader, binary.LittleEndian, &fullsize); err != nil {
		return
	}

	// 封包太大
	if self.maxPacketSize > 0 && int(fullsize) > self.maxPacketSize {
		err = ErrPackageTooBig
		return
	}

	// 序列号不匹配
	if self.recvser != ser {
		err = ErrPackageTagNotMatch
		return
	}

	dataSize := fullsize - PackageHeaderSize
	if dataSize < 0 {
		err = ErrPackageDataSizeInvalid
		return
	}

	// 读取数据
	msgBytes := make([]byte, dataSize)
	if _, err = io.ReadFull(self.conn, msgBytes); err != nil {
		return
	}

	data = msgBytes

	// 增加序列号值
	self.recvser++

	return
}

// 将一个封包发送到socket
func (self *TLVStream) Write(msgid uint32, data []byte) (err error) {

	// 防止将Send放在go内造成的多线程冲突问题
	self.sendtagGuard.Lock()
	defer self.sendtagGuard.Unlock()

	self.outputHeadBuffer.Reset()

	// 写ID
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, msgid); err != nil {
		return err
	}

	// 写序号
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, self.sendser); err != nil {
		return err
	}

	// 写包大小
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, uint16(len(data)+PackageHeaderSize)); err != nil {
		return err
	}

	// 发包头
	if err = self.writeFull(self.outputHeadBuffer.Bytes()); err != nil {
		return err
	}

	// 发包内容
	if err = self.writeFull(data); err != nil {
		return err
	}

	// 增加序号值
	self.sendser++

	return
}

// 完整发送所有封包
func (self *TLVStream) writeFull(p []byte) error {

	sizeToWrite := len(p)

	for {

		n, err := self.outputWriter.Write(p)

		if err != nil {
			return err
		}

		if n >= sizeToWrite {
			break
		}

		copy(p[0:sizeToWrite-n], p[n:sizeToWrite])
		sizeToWrite -= n
	}

	return nil

}

const sendTotalTryCount = 100

func (self *TLVStream) Flush() error {

	var err error
	for tryTimes := 0; tryTimes < sendTotalTryCount; tryTimes++ {

		err = self.outputWriter.Flush()

		// 如果没写完, flush底层会将没发完的buff准备好, 我们只需要重新调一次flush
		if err != io.ErrShortWrite {
			break
		}
	}

	return err
}

// 关闭
func (self *TLVStream) Close() error {
	return self.conn.Close()
}

// 裸socket操作
func (self *TLVStream) Raw() net.Conn {
	return self.conn
}

// 封包流 relay模式: 在封包头有clientid信息
func NewTLVStream(conn net.Conn) *TLVStream {

	s := &TLVStream{
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

func errToResult(err error) cellnet.Result {

	if err == nil {
		return cellnet.Result_OK
	}

	switch err {
	case ErrPackageTagNotMatch, ErrPackageDataSizeInvalid, ErrPackageTooBig:
		return cellnet.Result_PackageCrack
	}

	switch n := err.(type) {
	case net.Error:
		if n.Timeout() {
			return cellnet.Result_SocketTimeout
		}
	}

	return cellnet.Result_SocketError
}
