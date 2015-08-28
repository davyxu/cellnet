/*
	dispatcher包提供以注册+回调方式的消息处理方式, 封装消息解包, 打包的过程

*/
package dispatcher

import (
	"github.com/davyxu/cellnet"
)

type DataDispatcher struct {
	contextMap map[int][]func(cellnet.CellID, interface{})
}

func (self *DataDispatcher) RegisterCallback(id int, f func(cellnet.CellID, interface{})) {

	// 事件
	em, ok := self.contextMap[id]

	// 新建
	if !ok {

		em = make([]func(cellnet.CellID, interface{}), 0)

	}

	em = append(em, f)

	self.contextMap[id] = em
}

func (self *DataDispatcher) Exists(id int) bool {
	_, ok := self.contextMap[id]

	return ok
}

func (self *DataDispatcher) Call(src cellnet.CellID, id int, data interface{}) {

	if carr, ok := self.contextMap[id]; ok {

		for _, c := range carr {

			c(src, data)
		}
	}
}

func NewPacketDispatcher() *DataDispatcher {
	return &DataDispatcher{
		contextMap: make(map[int][]func(cellnet.CellID, interface{})),
	}

}
