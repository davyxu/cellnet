package socket

import (
	"fmt"

	"github.com/davyxu/cellnet"

	"github.com/golang/protobuf/proto"
)

type MessageLogInfo struct {
	Dir string
	ses cellnet.Session
	pkt *cellnet.Packet

	meta *cellnet.MessageMeta
}

func (self *MessageLogInfo) PeerName() string {
	return self.ses.FromPeer().Name()
}

func (self *MessageLogInfo) SessionID() int64 {
	return self.ses.ID()
}

func (self *MessageLogInfo) MsgName() string {

	if self.meta == nil {
		return ""
	}

	return self.meta.Name
}

func (self *MessageLogInfo) MsgID() uint32 {
	return self.pkt.MsgID
}

func (self *MessageLogInfo) MsgSize() int {
	return len(self.pkt.Data)
}

func (self *MessageLogInfo) MsgString() string {
	if self.meta == nil {
		return fmt.Sprintf("%v", self.pkt.Data)
	}

	rawMsg, err := cellnet.ParsePacket(self.pkt, self.meta.Type)
	if err != nil {
		return err.Error()
	}

	return rawMsg.(proto.Message).String()
}

// 是否启用消息日志
var EnableMessageLog bool = true

func msgLog(dir string, ses cellnet.Session, pkt *cellnet.Packet) {

	if !EnableMessageLog {
		return
	}

	info := &MessageLogInfo{
		Dir:  dir,
		ses:  ses,
		pkt:  pkt,
		meta: cellnet.MessageMetaByID(pkt.MsgID),
	}

	// 找到消息需要屏蔽
	if _, ok := msgMetaByID[info.MsgID()]; ok {
		return
	}

	if msgLogHook == nil || (msgLogHook != nil && msgLogHook(info)) {

		log.Debugf("#%s(%s) sid: %d %s size: %d | %s", info.Dir, info.PeerName(), info.SessionID(), info.MsgName(), info.MsgSize(), info.MsgString())

	}

}

var msgLogHook func(*MessageLogInfo) bool

func HookMessageLog(hook func(*MessageLogInfo) bool) {
	msgLogHook = hook
}

var msgMetaByID = make(map[uint32]*cellnet.MessageMeta)

func BlockMessageLog(msgName string) {
	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		log.Errorf("msg log block not found: %s", msgName)
		return
	}

	msgMetaByID[meta.ID] = meta

}
