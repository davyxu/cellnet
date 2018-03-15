package proc

import (
	"github.com/davyxu/cellnet"
)

type ProcessorBundle interface {
	SetEventTransmitter(v cellnet.MessageTransmitter)

	SetEventHooker(v cellnet.EventHooker)

	SetEventCallback(v cellnet.EventCallback)
}

type ProcessorBinder func(bundle ProcessorBundle, userCallback cellnet.EventCallback)

var (
	procByName = map[string]ProcessorBinder{}
)

// 注册事件处理器，内部及自定义收发流程时使用
func RegisterEventProcessor(procName string, f ProcessorBinder) {

	procByName[procName] = f
}

// 绑定固定回调处理器
func BindProcessorHandler(peer cellnet.Peer, procName string, userCallback cellnet.EventCallback) {

	if proc, ok := procByName[procName]; ok {

		bundle := peer.(ProcessorBundle)

		proc(bundle, userCallback)

	} else {
		panic("processor not found:" + procName)
	}
}
