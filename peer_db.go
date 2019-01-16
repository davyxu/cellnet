package cellnet

type RedisPoolOperator interface {
	// 获取
	Operate(callback func(rawClient interface{}) interface{}) interface{}
}

type RedisConnector interface {
	GenericPeer

	SetPassword(v string)

	SetDBIndex(v int)

	SetConnectionCount(v int)
}

type MySQLConnector interface {
	GenericPeer

	SetPassword(v string)

	SetConnectionCount(v int)
}
