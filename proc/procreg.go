package proc

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"sort"
)

type ProcessorBinder func(bundle ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{})

var (
	procByName = map[string]ProcessorBinder{}
)

// 注册事件处理器，内部及自定义收发流程时使用
func RegisterProcessor(procName string, f ProcessorBinder) {

	if _, ok := procByName[procName]; ok {
		panic("duplicate peer type: " + procName)
	}

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

func getPackageByCodecName(name string) string {
	switch name {
	case "gorillaws.ltv":
		return "github.com/davyxu/cellnet/proc/gorillaws"
	case "http":
		return "github.com/davyxu/cellnet/proc/http"
	case "tcp.ltv":
		return "github.com/davyxu/cellnet/proc/tcp"
	case "udp.ltv":
		return "github.com/davyxu/cellnet/proc/udp"
	default:
		return "package/to/your/proc"
	}
}

// 绑定固定回调处理器, procName来源于RegisterProcessor注册的处理器，形如: 'tcp.ltv'
func BindProcessorHandler(peer cellnet.Peer, procName string, userCallback cellnet.EventCallback, args ...interface{}) {

	if proc, ok := procByName[procName]; ok {

		bundle := peer.(ProcessorBundle)

		proc(bundle, userCallback, args...)

	} else {
		panic(fmt.Sprintf("processor not found '%s'\ntry to add code below:\nimport (\n  _ \"%s\"\n)\n\n",
			procName,
			getPackageByCodecName(procName)))
	}
}
