package cellnet

import (
	"html/template"
	"time"
)

// 基本的通用Peer
type GenericPeer interface {
	Peer
	PeerProperty
}

type HTTPAcceptor interface {
	GenericPeer

	// 设置http文件服务虚拟地址和文件系统根目录
	SetFileServe(dir string, root string)

	// 设置模板文件地址
	SetTemplateDir(dir string)

	// 设置http模板的分隔符，解决默认{{ }}冲突问题
	SetTemplateDelims(delimsLeft, delimsRight string)

	// 设置模板的扩展名，默认: .tpl .html
	SetTemplateExtensions(exts []string)

	// 设置模板函数入口
	SetTemplateFunc(f []template.FuncMap)
}

// HTTP连接器接口
type HTTPConnector interface {
	GenericPeer
	Request(method string, raw interface{}) (interface{}, error)
}

type SessionManagerSetter interface {

	// 将peer.SesssionManager设置给Connector，Connector的session就分配到一个SessionManager中
	SetSessionManager(sessionManager interface{})
}

// TCP连接器
type TCPConnector interface {
	GenericPeer

	SessionManagerSetter

	// 设置重连时间
	SetReconnectDuration(time.Duration)

	// 获取重连时间
	ReconnectDuration() time.Duration

	// 默认会话
	Session() Session
}

// TCP接受器，具备会话访问
type TCPAcceptor interface {
	GenericPeer

	SessionManagerSetter

	// 访问会话
	SessionAccessor
}

// UDP连接器
type UDPConnector interface {
	GenericPeer

	// 默认会话
	Session() Session
}
