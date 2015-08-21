package dispatcher

import (
	"github.com/davyxu/cellnet"
)

type PacketDispatcher struct {
	contextMap map[int][]func(cellnet.CellID, *cellnet.Packet)
}

func (self *PacketDispatcher) RegisterCallback(id int, f func(cellnet.CellID, *cellnet.Packet)) {

	// 事件
	em, ok := self.contextMap[id]

	// 新建
	if !ok {

		em = make([]func(cellnet.CellID, *cellnet.Packet), 0)

	}

	em = append(em, f)

	self.contextMap[id] = em
}

func (self *PacketDispatcher) Exists(id int) bool {
	_, ok := self.contextMap[id]

	return ok
}

func (self *PacketDispatcher) Call(src cellnet.CellID, pkt *cellnet.Packet) {

	if carr, ok := self.contextMap[int(pkt.MsgID)]; ok {

		for _, c := range carr {

			c(src, pkt)
		}
	}
}

func NewPacketDispatcher() *PacketDispatcher {
	return &PacketDispatcher{
		contextMap: make(map[int][]func(cellnet.CellID, *cellnet.Packet)),
	}

}
