# cellnet
简单,方便,高效的跨平台服务器网络库


# 特性

## 异步单线程多进程架构
  
* 无需处理繁琐的多线程安全问题

* 底层IO仍然使用goroutine进行处理, 保证IO吞吐率

* 发送时自动合并封包(性能效果决定于实际请求和发送比例)

## 数据协议

* 封包类型采用Type-Length-Value的私有tcp封包, 自带序列号防御简单的封包复制

* 消息统一使用Protobuf格式进行通信, 有自定义需求可以修改实现

* 自动生成消息ID

## RPC

* 异步远程过程调用

## 日志
* 分级日志

* 可以方便的通过日志查看收发消息(Protobuf)的每一个字段消息

# 第三方库依赖

* github.com/golang/protobuf/proto

* github.com/davyxu/golog

# 性能测试

命令行: go test -v github.com/davyxu/cellnet/benchmark/io

平台: Windows 7 x64/CentOS 6.5 x64

测试用例: localhost 1000连接 同时对服务器进行实时PingPong测试

配置1: i7 6700 3.4GHz 8核

QPS: 12.5w

配置2: i5 4590 3.3GHz 4核

QPS: 10.1w

# 例子
## Echo
```go


func server() {

	queue := cellnet.NewEventQueue()

	evd := socket.NewAcceptor(queue).Start("127.0.0.1:7201")

	socket.RegisterMessage(evd, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.String())

		ses.Send(&gamedef.TestEchoACK{
			Content: msg.String(),
		})

	})

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	evd := socket.NewConnector(queue).Start("127.0.0.1:7201")

	socket.RegisterMessage(evd, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.String())

		signal.Done(1)
	})

	socket.RegisterMessage(evd, "gamedef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		ses.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

	})

	queue.StartLoop()
}

```

# 文件夹功能

benchmark\		性能测试用例

db\				db封装

proto\			cellnet内部的proto

protoc-gen-msg\ 消息id生成

rpc\			异步远程过程调用封装

socket\			套接字,拆包等封装

sample\			测试用例/例子
   
	close\		发送消息后保证消息送达后再断开连接
	
   	echo\		常见的pingpong测试， 最简单的例子
	
	mgo\		mongodb异步读取例子
	
	rpc\		异步远程过程调用
	
	timer\		异步计时器
	
	
util\			工具库

# FAQ

问: 为什么接收消息回调必须需要手动转换类型, 例如: msg := content.(*gamedef.TestEchoACK), 而不是参数上写成参数类型?
	
答: cellnet这么设计是考虑到可以将参数进行多层传递, 如果写成不同消息类型, 传递就麻烦很多

这里鼓励消息注册者可以进行消息注册函数的封装, 在网关里, 就对socket.RegisterSessionMessage进行封装
	
```golang
func RegisterMessage(msgName string, userHandler func(interface{}, cellnet.Session, int64)) {

	msgMeta := cellnet.MessageMetaByName(msgName)

	if msgMeta == nil {
		log.Errorf("message register failed, %s", msgName)
		return
	}

	for _, conn := range routerConnArray {

		conn.AddCallback(msgMeta.ID, func(data interface{}) {

			if ev, ok := data.(*relayEvent); ok {

				rawMsg, err := cellnet.ParsePacket(ev.Packet, msgMeta.Type)

				if err != nil {
					log.Errorln("unmarshaling error:\n", err)
					return
				}

				msgContent := rawMsg.(interface {
					String() string
				}).String()				

				userHandler(rawMsg, ev.Ses, ev.ClientID)

			}

		})
	}

}

```

再来一个外层封装
```golang
func RegisterExternMessage(msgName string, userHandler func(interface{}, *Player)) {

	backend.RegisterMessage(msgName, func(content interface{}, routerSes cellnet.Session, clientid int64) {

		if player, ok := PlayerByID[clientid]; ok {

			userHandler(content, player)
		}
	})

}
```
# 版本历史
2017.1  v2版本 [详细请查看](https://github.com/davyxu/cellnet/blob/master/CHANGES.md)

2015.8	v1版本

# Wiki
https://github.com/davyxu/cellnet/wiki

这里有文档和架构,设计解析


# 贡献者

欢迎提供dev分支的pull request

bug请直接通过issue提交

凡提交代码和建议, bug的童鞋, 均会在下列贡献者名单者出现

viwii(viwii@sina.cn), 提供一个可能造成死锁的bug

IronsDu(duzhongwei@qq.com), 大幅度性能优化

Chris Lonng(chris@lonng.org), 提供一个最大封包约束造成服务器间连接断开的bug

# 备注

感觉不错请star, 谢谢!

博客: http://www.cppblog.com/sunicdavy

知乎: http://www.zhihu.com/people/sunicdavy

提交bug及特性: https://github.com/davyxu/cellnet/issues

贡献代码: https://github.com/davyxu/cellnet/pulls