package cellnet

type Peer interface {

	// 开启
	Start(address string) Peer

	// 关闭
	Stop()

	// 名字
	SetName(string)
	Name() string

	// 获取一个连接
	Get(uint32) Session

	// 广播
	Broardcast(interface{})

	// 遍历连接
	Iterate(func(Session) bool)

	// 连接数量
	Count() int
}
