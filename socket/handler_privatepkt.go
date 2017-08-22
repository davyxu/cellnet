package socket

import (
	"bytes"
	"encoding/binary"
	"github.com/davyxu/cellnet"
	"io"
	"sync"
)

type PrivatePacketReader struct {
	recvser uint16
}

func (self *PrivatePacketReader) Call(ev *cellnet.Event) {

	headReader := bytes.NewReader(ev.Data)

	// 读取序号
	var ser uint16
	if err := binary.Read(headReader, binary.LittleEndian, &ser); err != nil {
		ev.SetResult(cellnet.Result_PackageCrack)
		return
	}

	// 读取ID
	if err := binary.Read(headReader, binary.LittleEndian, &ev.MsgID); err != nil {
		ev.SetResult(cellnet.Result_PackageCrack)
		return
	}

	// 读取Payload大小
	var bodySize uint32
	if err := binary.Read(headReader, binary.LittleEndian, &bodySize); err != nil {
		ev.SetResult(cellnet.Result_PackageCrack)
		return
	}

	maxPacketSize := ev.Ses.FromPeer().(SocketOptions).MaxPacketSize()
	// 封包太大
	if maxPacketSize > 0 && int(bodySize) > maxPacketSize {
		ev.SetResult(cellnet.Result_PackageCrack)
		return
	}

	// 序列号不匹配
	if self.recvser != ser {
		ev.SetResult(cellnet.Result_PackageCrack)
		return
	}

	reader := ev.Ses.(interface {
		DataSource() io.ReadWriter
	}).DataSource()

	// 读取数据
	dataBuffer := make([]byte, bodySize)
	if _, err := io.ReadFull(reader, dataBuffer); err != nil {
		ev.SetResult(cellnet.Result_PackageCrack)
		return
	}

	ev.Data = dataBuffer

	// 增加序列号值
	self.recvser++
}

func NewPrivatePacketReader() cellnet.EventHandler {
	return &PrivatePacketReader{
		recvser: 1,
	}
}

type PrivatePacketWriter struct {
	sendser      uint16
	sendtagGuard sync.RWMutex
}

func (self *PrivatePacketWriter) Call(ev *cellnet.Event) {

	// 防止将Send放在go内造成的多线程冲突问题
	self.sendtagGuard.Lock()
	defer self.sendtagGuard.Unlock()

	var outputHeadBuffer bytes.Buffer

	// 写序号
	if err := binary.Write(&outputHeadBuffer, binary.LittleEndian, self.sendser); err != nil {
		ev.SetResult(cellnet.Result_PackageCrack)
		return
	}

	// 写ID
	if err := binary.Write(&outputHeadBuffer, binary.LittleEndian, ev.MsgID); err != nil {
		ev.SetResult(cellnet.Result_PackageCrack)
		return
	}

	// 写包大小
	if err := binary.Write(&outputHeadBuffer, binary.LittleEndian, uint32(len(ev.Data))); err != nil {
		ev.SetResult(cellnet.Result_PackageCrack)
		return
	}

	binary.Write(&outputHeadBuffer, binary.LittleEndian, ev.Data)

	// 增加序号值
	self.sendser++

	ev.Data = outputHeadBuffer.Bytes()
}

func NewPrivatePacketWriter() cellnet.EventHandler {
	return &PrivatePacketWriter{
		sendser: 1,
	}
}
