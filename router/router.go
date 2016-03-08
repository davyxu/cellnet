package router

import (
	"sync"

	"github.com/davyxu/cellnet"
)

type routeInfo struct {
	name string
	ses  cellnet.Session
}

var (
	routerMap      = make(map[uint32]*routeInfo)
	routerMapGuard sync.RWMutex
)

// 注册一个后台服务器的连接
func registerBackend(ses cellnet.Session, name string) {

	routerMapGuard.Lock()

	// 更新路由表中的连接
	for _, ri := range routerMap {

		if ri.name == name {
			ri.ses = ses
		}

	}

	routerMapGuard.Unlock()
}

// 后台服务器断开连接时, 及时清理连接
func closeBackend(ses cellnet.Session) {
	routerMapGuard.Lock()

	// 更新路由表中的连接
	for _, ri := range routerMap {

		if ri.ses == ses {
			ri.ses = nil
		}

	}

	routerMapGuard.Unlock()

}

// 根据消息id决定路由目的地
func getRelaySession(msgid uint32) cellnet.Session {

	routerMapGuard.RLock()

	defer routerMapGuard.RUnlock()

	if ri, ok := routerMap[msgid]; ok {

		return ri.ses
	}

	return nil

}

type RelayMethod int

const (
	// 广播到后台所有服务器
	RelayMethod_BroardcastToAllBackend RelayMethod = iota

	// 按照白名单准确投递
	RelayMethod_WhiteList
)

var relayMethod RelayMethod = RelayMethod_BroardcastToAllBackend

// 设置路由模式
func SetRelayMethod(method RelayMethod) {
	relayMethod = method
}

// 注册消息路由方法
func RelayMessage(targetSvcName string, messageName string) {

	meta := cellnet.MessageMetaByName(messageName)
	if meta == nil {
		log.Errorf("relay message not found: %s, target: %s", messageName, targetSvcName)
		return
	}

	routerMap[meta.ID] = &routeInfo{name: targetSvcName}
}
