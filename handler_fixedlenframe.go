package cellnet

import (
	"io"
	"sync"
)

type FixedLengthFrameReader struct {
	headerBuffer []byte
}

func (self *FixedLengthFrameReader) Call(ev *Event) {

	reader := ev.Ses.(interface {
		DataSource() io.ReadWriter
	}).DataSource()

	_, err := io.ReadFull(reader, self.headerBuffer)

	if err != nil {
		ev.SetResult(Result_SocketError)
		return
	}

	ev.Data = self.headerBuffer
}

func NewFixedLengthFrameReader(size int) EventHandler {
	return &FixedLengthFrameReader{
		headerBuffer: make([]byte, size),
	}
}

type FixedLengthFrameWriter struct {
	sendser      uint16
	sendtagGuard sync.RWMutex
}

func (self *FixedLengthFrameWriter) Call(ev *Event) {

	writer := ev.Ses.(interface {
		DataSource() io.ReadWriter
	}).DataSource()

	err := writeFull(writer, ev.Data)

	if err != nil {
		ev.SetResult(Result_PackageCrack)
		return
	}

}

// 完整发送所有封包
func writeFull(writer io.ReadWriter, p []byte) error {

	sizeToWrite := len(p)

	for {

		n, err := writer.Write(p)

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

func NewFixedLengthFrameWriter() EventHandler {
	return &FixedLengthFrameWriter{}
}
