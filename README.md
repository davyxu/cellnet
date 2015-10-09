# Cellnet
简单,方便,高效的Go语言的游戏服务器框架


# 特性
## 异步单线程多进程架构
  
* 无需处理繁琐的多线程安全问题
* 底层IO仍然使用goroutine进行处理, 保证IO吞吐率
* 性能敏感的业务拆离为单独进程进行处理

## 数据协议
* 封包类型采用Type-Length-Value的私有tcp封包, 自带序列号防御简单的封包复制
* 消息统一使用Protobuf格式进行通信

## 网关
* 基本的网关透传框架
* 广播,列表广播


## RPC
* 异步远程过程调用

## 异步MongoDB
* 提供KV数据库的基本抽象

## 模块化
* 鼓励使用统一的模块化命名及拆分方法进行隔离降偶


# 第三方库依赖

* github.com/golang/protobuf/proto
* github.com/BurntSushi/toml
* gopkg.in/mgo.v2


# 例子
## Echo
```go


func server() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewAcceptor(pipe).Start("127.0.0.1:7234")

	socket.RegisterSessionMessage(evq, coredef.TestEchoACK{}, func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("server recv:", msg.String())

		ses.Send(&coredef.TestEchoACK{
			Content: proto.String(msg.String()),
		})

	})

	pipe.Start()

}

func client() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewConnector(pipe).Start("127.0.0.1:7234")

	socket.RegisterSessionMessage(evq, coredef.TestEchoACK{}, func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

	})

	socket.RegisterSessionMessage(evq, coredef.SessionConnected{}, func(content interface{}, ses cellnet.Session) {

		ses.Send(&coredef.TestEchoACK{
			Content: proto.String("hello"),
		})

	})

	pipe.Start()
}

```

# 新品预告!

## 日志
* 支持分级

## 消息日志
* 可以方便的通过日志查看收发消息(Protobuf)的每一个字段消息

## MongoDB
* DB内存映射框架
* DB存储日志


## 网关
* 可定制的消息路由规则

# 备注
本人开发环境是Windows, 所以工具、Shell相关都会以Windows为主
后期会对Linux平台的开发环境进行支持

推荐使用LiteIDE(https://github.com/visualfc/liteide)进行开发及查看


博客: http://www.cppblog.com/sunicdavy

知乎: http://www.zhihu.com/people/xu-bo-62-87

技术讨论组: 309800774 加群请说明cellnet

邮箱: sunicdavy@qq.com
