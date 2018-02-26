# cellnet
 [![Build Status][3]][4] [![Go Report Card][5]][6] [![MIT licensed][11]][12] [![GoDoc][1]][2]

[1]: https://godoc.org/github.com/davyxu/cellnet?status.svg
[2]: https://godoc.org/github.com/davyxu/cellnet
[3]: https://travis-ci.org/davyxu/cellnet.svg?branch=v4
[4]: https://travis-ci.org/davyxu/cellnet
[5]: https://goreportcard.com/badge/github.com/davyxu/cellnet
[6]: https://goreportcard.com/report/github.com/davyxu/cellnet
[11]: https://img.shields.io/badge/license-MIT-blue.svg
[12]: LICENSE

cellnet是一个组件化、高扩展性、高性能的开源服务器网络库


# 特性

## 传输协议支持
- TCP

    TCP连接器的重连，侦听器的优雅重启。

- UDP

    纯UDP裸包收发

- HTTP

    支持json及form的收发及封装

- KCP

    支持KCP over UDP

## 混合编码（Codec）支持

* cellnet内建支持以下数据编码:
    - Google Protobuf (https://github.com/google/protobuf)

    - 云风的sproto (https://github.com/cloudwu/sproto)

        能方便lua的处理, 本身结构比protobuf解析更简单

    - json
        适合与第三方服务器通信

    - 二进制协议(https://github.com/davyxu/goobjfmt)

       内存流直接序列化, 适用于服务器内网传输

    可以通过codec包自行添加新的编码格式

* 支持混合编码收发

    优势：无需改动代码，只需调整消息注册方式，即可达成收发不同编码的封包

    - 与Unity3D+Lua使用sproto通信

    - 与其他语言编写的服务器使用protobuf

    - 与web服务器使用json通信

## 队列及IO
  
* 支持多个队列, 实现单线程/多线程收发处理消息

* 发送时自动合并封包(性能效果决定于实际请求和发送比例)

## RPC

* 异步/同步远程过程调用

## 消息日志
* 可以方便的通过日志查看收发消息的每一个字段消息

# 第三方库依赖

* github.com/davyxu/golog

* github.com/davyxu/goobjfmt

# 编码包可选支持

* github.com/golang/protobuf

* github.com/davyxu/gosproto

# 获取+编译

```
	go get -u -v github.com/davyxu/cellnet

```


# 样例
```golang


func server() {

	queue := cellnet.NewEventQueue()

	p := socket.NewAcceptor(queue).Start("127.0.0.1:7201")

	cellnet.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.Content)

		ev.Send(&gamedef.TestEchoACK{
			Content: msg.String(),
		})

	})

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	p := socket.NewConnector(queue).Start("127.0.0.1:7301")

	cellnet.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.Content)
	})

	cellnet.RegisterMessage(p, "coredef.SessionConnected", func(ev *cellnet.Event) {

		log.Debugln("client connected")

		ev.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

	})

	cellnet.RegisterMessage(p, "coredef.SessionConnectFailed", func(ev *cellnet.Event) {

		msg := ev.Msg.(*coredef.SessionConnectFailed)

		log.Debugln(msg.Reason)

	})

	queue.StartLoop()
}

```

# 目录功能

```
benchmark           性能测试用例

codec               编码支持，以及编码注册

    binary          二进制格式编码(github.com/davyxu/goobjfmt)

    httpform        http表单格式

    json            json编码格式

examples            例子

peer                各种协议的端实现，以及端注册入口及复用组件

    http            HTTP协议处理流程及端封装

    tcp             TCP协议处理流程及端封装

    udp             UDP协议处理流程及端封装

proc                各种处理器实现，以及处理器注册入口

    http            HTTP消息处理及文件服务实现

    kcp             在UDP peer上构建的KCP协议支持

    rpc             远程过程调用支持

    tcp             在TCP peer上构建的tcp处理器集合

    udp             在UDP peer上构建的udp处理器集合

tests               测试用例

timer               计时器接口

util                工具库

```


# 运行聊天例子

## 运行 服务器

```bash
cd examples/chat/server

go run main.go
```

## 运行 客户端

```bash
cd examples/chat/client

go run main.go
```

随后, 在命令行中输入hello后打回车, 就可以看到服务器返回

```

sid1 say: hello

```

# 术语及概念

## 队列

队列在cellnet中使用cellnet.Queue接口, 底层由带缓冲的channel实现

在cellnet中, 队列根据实际逻辑需要定制数量. 但一般情况下, 仅需要1个队列

使用下列代码使用队列

```golang
    queue := cellnet.NewEventQueue()

    // 启动队列
    queue.StartLoop()

    // 这里添加队列使用代码

    // 等待队列结束, 调用queue.StopLoop(0)将退出阻塞
    queue.Wait()
```


当在多线程环境中, 需要将逻辑排队执行时, 我们只需要这样写:

```golang
    queue.Post(func() {
		fmt.Println("hello")
	})

```

我们使用的消息, 在底层就是通过queue.Post() 传入我们给定的回调进行处理的

### 提示

队列对于使用cellnet的服务器程序不是必须的.不使用队列时, 所有消息的处理将是并发在多线程下


## 侦听和接受连接

使用如下代码创建一个接受器(Acceptor)

```golang
    queue := cellnet.NewEventQueue()

    peer := socket.NewAcceptor(queue)

    peer.Start("127.0.0.1:8801")
```

底层将自动完成侦听, 接受连接, 错误抛出, 断开处理等复杂操作

### 提示

- cellnet将系统事件, 错误和用户消息均视为消息, 并可以被注册回调后处理, 接口都是统一的


接受器可接收的系统消息可以用下面代码注册并响应:

```golang
cellnet.RegisterMessage(peer, "coredef.SessionAccepted", func(ev *cellnet.Event) {		
    // 其他会话连接时
})

cellnet.RegisterMessage(peer, "coredef.SessionAcceptFailed", func(ev *cellnet.Event) {		
    // 其他会话连接失败时
})
```

## 发起连接

使用如下代码创建一个连接器(Connector)

```golang
    queue := cellnet.NewEventQueue()

    peer := socket.NewConnector(queue)

    peer.Start("127.0.0.1:8801")
```

底层将自动完成连接; 如果发生断开, 可以通过如下代码设置自动重连

```golang
    // 设置连接超时2秒后自动重连
    peer.(socket.Connector).SetAutoReconnectSec(2)
```

连接器也可以接收系统事件, 如:
```golang
cellnet.RegisterMessage(peer, "coredef.SessionConnectFailed", func(ev *cellnet.Event) {		
    // 会话连接失败
})

```

### 提示:

- 端
    
    cellnet中, Connector和Acceptor被统称为Peer, 即"端", 当连接器和接受器建立连接后, 两个"端"的概念, 接口和使用均是相同的

# 会话连接(Session)

建立连接后, 这个连接在cellnet中称为Session

## Session ID
每个会话拥有一个64位ID, cellnet底层保证在一个Peer中不会重复

可以通过Session.ID() 获得ID

## 获得Session
Session可以通过以下途径获得:

- 通过cellnet.RegisterMessage注册回调后, 通过回调参数*cellnet.Event中的Ses获得

- 如果Peer是Connector, 可以通过如下代码获得连接器上默认连接

```golang
    ses := peer.(socket.Connector).DefaultSession()
```

- Peer可以通过SessionAccessor接口以多种方式获得Session, 如:
```golang

    // 通过Session.ID获得
    GetSession(int64) Session

    // 遍历这个Peer的所有Session
    VisitSession(func(Session) bool)

```

## 接收消息

消息通过cellnet.Event传递

如需获得消息, 我们使用如下代码获得消息
    
```golang
    msg := ev.Msg.(*MsgPackage.YourMsgType)
```

## 接收系统事件

如需接收Session连接断开事件, 使用如下代码
```golang
cellnet.RegisterMessage(peer, "coredef.SessionClosed", func(ev *cellnet.Event) {		
    // 会话断开时
})

```


## 会话发送消息

一般情况下, 我们的消息使用结构体实现. 使用protobuf工具链, sproto工具链可以直接生成这些结构体.

通过Session.Send()发送一个结构体指针, 如:
```golang
    ses.Send(&chatproto.ChatREQ{
			Content: str,
		})
```




### 提示

- 使用Event.Send方式回消息

```golang
    cellnet.RegisterMessage(peer, "chatproto.ChatACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*chatproto.ChatACK)

        // 使用这种方法回应消息与rpc系统统一, 便于底层优化
		ev.Send( msg )
	})

```

- 不要缓存Event

Event是消息处理的上下文, 不建议缓存Event





# FAQ

* 这个代码的入口在哪里? 怎么编译为exe?

    本代码是一个网络库, 需要根据需求, 整合逻辑

* 混合编码有何用途?

    在与多种语言写成的服务器进行通信时, 可以使用不同的编码,
    最终在逻辑层都是统一的结构能让逻辑编写更加方便, 无需关注底层处理细节

* 内建支持的二进制协议能与其他语言写成的网络库互通么?

    完全支持, 但内建二进制协议支持更适合网关与后台服务器.
    不建议与客户端通信中使用, 二进制协议不会忽略使用默认值的字段

* 所有的例子都是单线程的, 能编写多线程的逻辑么?

    完全可以, cellnet并没有全局的队列, 只需在Acceptor和Connector创建时,
    传入不同的队列, socket收到的消息就会被放到这个队列中
    传入空队列时, 使用并发方式(io线程)调用处理回调

* cellnet有网关和db支持么?

    cellnet专注于服务器底层.你可以根据自己需要编写网关及db支持

* 怎样定制私有TCP封包?

* 哪里有cellnet的完整例子?

    CellOrigin是基于cellnet开发的一套Unity3D客户端服务器框架
    https://github.com/davyxu/cellorigin

# 版本历史
2018.8  v4版本 [详细请查看](https://github.com/davyxu/cellnet/blob/v4/CHANGES.md)

2017.8  v3版本 [详细请查看](https://github.com/davyxu/cellnet/blob/v3/CHANGES.md)

2017.1  v2版本 [详细请查看](https://github.com/davyxu/cellnet/blob/v2/CHANGES.md)

2015.8	v1版本


# 贡献者

按贡献时间排序，越靠前表示越新的贡献

superikw(https://github.com/superikw), 测试出一个websocket接口并发发送问题

bruce.hu(https://github.com/hxdhero), 测试出一个竞态冲突的bug

M4tou(https://github.com/mutousay), 协助解决RPC异步超时回调处理

chuan.li(https://github.com/blade-226), 提供一个没有在io线程编码的bug

Chris Lonng(https://github.com/lonnng), 提供一个最大封包约束造成服务器间连接断开的bug

IronsDu(https://github.com/IronsDu), 大幅度性能优化

viwii(viwii@sina.cn), 提供一个可能造成死锁的bug


# 备注

感觉不错请star, 谢谢!

博客: http://www.cppblog.com/sunicdavy

知乎: http://www.zhihu.com/people/sunicdavy

提交bug及特性: https://github.com/davyxu/cellnet/issues
