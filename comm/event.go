package comm

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
)

// 接收到封包
type RecvPacketEvent struct {
	Ses    cellnet.Session
	Reader *util.BinaryReader
}
