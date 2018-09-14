package proc

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"sort"
)

type ProcessorBinder func(bundle ProcessorBundle, userCallback cellnet.EventCallback)

var (
	procByName = map[string]ProcessorBinder{}
)

// 注册事件处理器，内部及自定义收发流程时使用
func RegisterProcessor(procName string, f ProcessorBinder) {

	procByName[procName] = f
}

// 获取处理器列表
func ProcessorList() (ret []string) {

	for name := range procByName {
		ret = append(ret, name)
	}

	sort.Strings(ret)
	return
}

// 绑定固定回调处理器, procName来源于RegisterProcessor注册的处理器，形如: 'tcp.ltv'
func BindProcessorHandler(peer cellnet.Peer, procName string, userCallback cellnet.EventCallback) {

	if proc, ok := procByName[procName]; ok {

		bundle := peer.(ProcessorBundle)

		proc(bundle, userCallback)

	} else {
		panic(fmt.Sprintf("processor not found, name: '%s'", procName))
	}
}
