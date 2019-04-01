package cellnet

import "time"

type RedisPoolOperator interface {
	// 获取
	Operate(callback func(rawClient interface{}) interface{}) interface{}
}

type RedisConnector interface {
	GenericPeer

	// 设置密码
	SetPassword(v string)

	// 设置连接数
	SetConnectionCount(v int)

	// 设置库索引
	SetDBIndex(v int)
}

type MySQLOperator interface {
	Operate(callback func(rawClient interface{}) interface{}) interface{}
}

type MySQLConnector interface {
	GenericPeer

	// 设置密码
	SetPassword(v string)

	// 设置连接数
	SetConnectionCount(v int)

	// 设置自动重连间隔， 0为默认值，关闭自动重连
	SetReconnectDuration(v time.Duration)

	ReconnectDuration() time.Duration
}
