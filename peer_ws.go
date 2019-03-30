package cellnet

import "time"

// Websocket接受器，具备会话访问
type WSAcceptor interface {
	GenericPeer

	// 访问会话
	SessionAccessor

	SetHttps(certfile, keyfile string)

	// 设置升级器
	SetUpgrader(upgrader interface{})

	// 查看当前侦听端口，使用host:0 作为Address时，socket底层自动分配侦听端口
	Port() int
}

// Websocket连接器
type WSConnector interface {
	GenericPeer

	// 设置重连时间
	SetReconnectDuration(time.Duration)

	// 获取重连时间
	ReconnectDuration() time.Duration

	// 默认会话
	Session() Session

	// 设置会话管理器 实现peer.SessionManager接口
	SetSessionManager(raw interface{})

	// 查看当前连接使用的端口
	Port() int
}
