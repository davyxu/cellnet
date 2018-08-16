package cellnet

import (
	"html/template"
	"net/http"
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
	Request(method, path string, raw interface{}) (interface{}, error)
}

type HTTPSession interface {
	Request() *http.Request
}

// TCP连接器
type TCPConnector interface {
	GenericPeer

	TCPSocketOption

	// 设置重连时间
	SetReconnectDuration(time.Duration)

	// 获取重连时间
	ReconnectDuration() time.Duration

	// 默认会话
	Session() Session

	// 设置会话管理器 实现peer.SessionManager接口
	SetSessionManager(raw interface{})

	Port() int
}

// TCP接受器，具备会话访问
type TCPAcceptor interface {
	GenericPeer

	// 访问会话
	SessionAccessor

	TCPSocketOption

	// 查看当前侦听端口，使用host:0 作为Address时，socket底层自动分配侦听端口
	Port() int
}

// UDP连接器
type UDPConnector interface {
	GenericPeer

	// 默认会话
	Session() Session
}

// UDP接受器
type UDPAcceptor interface {

	// 底层使用TTL做session生命期管理，超时时间越短，内存占用越低
	SetSessionTTL(dur time.Duration)
}

// Websocket接受器，具备会话访问
type WSAcceptor interface {
	GenericPeer

	SetHttps(certfile, keyfile string)

	// 访问会话
	SessionAccessor
}

// TCP
type TCPSocketOption interface {
	// 收发缓冲大小，默认-1
	SetSocketBuffer(readBufferSize, writeBufferSize int, noDelay bool)

	// 设置最大的封包大小
	SetMaxPacketSize(maxSize int)

	// 设置读写超时，默认0，不超时
	SetSocketDeadline(read, write time.Duration)
}
