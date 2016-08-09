# Cellnet
简单,方便,高效的Go语言的游戏服务器底层


# 特性
## 异步单线程多进程架构
  
* 无需处理繁琐的多线程安全问题
* 底层IO仍然使用goroutine进行处理, 保证IO吞吐率
* 性能敏感的业务拆离为单独进程进行处理

## 数据协议
* 封包类型采用Type-Length-Value的私有tcp封包, 自带序列号防御简单的封包复制
* 消息统一使用Protobuf格式进行通信

## RPC
* 异步远程过程调用

## 模块化
* 鼓励使用统一的模块化命名及拆分方法进行隔离降偶

## 日志
* 分级日志
* 可以方便的通过日志查看收发消息(Protobuf)的每一个字段消息

# 第三方库依赖

* github.com/golang/protobuf/proto
* github.com/davyxu/golog
* gopkg.in/mgo.v2


# 例子
## Echo
```go


func server() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewAcceptor(pipe).Start("127.0.0.1:7234")

	socket.RegisterSessionMessage(evq, "coredef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.TestEchoACK)

		log.Debugln("server recv:", msg.String())

		ses.Send(&coredef.TestEchoACK{
			Content: msg.String,
		})

	})

	pipe.Start()

}

func client() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewConnector(pipe).Start("127.0.0.1:7234")

	socket.RegisterSessionMessage(evq, "coredef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.TestEchoACK)

		log.Debugln("client recv:", msg.String())

	})

	socket.RegisterSessionMessage(evq, "coredef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		ses.Send(&coredef.TestEchoACK{
			Content: "hello",
		})

	})

	pipe.Start()
}

```

# TODO


## MongoDB

* DB存储日志


# Wiki
https://github.com/davyxu/cellnet/wiki
这里有文档和架构,设计解析


# 备注

感觉不错请star, 谢谢!

博客: http://www.cppblog.com/sunicdavy

知乎: http://www.zhihu.com/people/xu-bo-62-87

邮箱: sunicdavy@qq.com

战魂小筑技术讨论群: 309800774 加群请说明cellnet

cellnet发问请直接@成都_黑色灵猫

# 贡献者
viwii
