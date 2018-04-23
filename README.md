# cellnet

![cellnetlogo](doc/logo.png)

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



# 版本历史
2018.5  v4版本 [详细请查看](https://github.com/davyxu/cellnet/blob/v4/CHANGES.md)

2017.8  v3版本 [详细请查看](https://github.com/davyxu/cellnet/blob/v3/CHANGES.md)

2017.1  v2版本 [详细请查看](https://github.com/davyxu/cellnet/blob/v2/CHANGES.md)

2015.8	v1版本

# 使用领域

cellnet经过多个版本的迭代，无论是作为初学者学习的范例，还是作为私用、商用项目的基础构建乃至核心技术层已经在业内广受了解及使用。

主要使用领域：

- 游戏服务器

- ARM设备网络通讯

- 证券软件内部RPC

# 特性 (Features)

## 传输协议支持
- TCP

    TCP连接器的重连，侦听器的优雅重启。

- UDP

    纯UDP裸包收发

- HTTP

    侦听器的优雅重启, 支持json及form的收发及封装。编写游戏服务器时，不再需要配合第三方HTTP服务器做后台。

## 编码(Codec)

* cellnet内建支持以下数据编码:
    - Google Protobuf (https://github.com/google/protobuf)

    - json
        适合与第三方服务器通信

    - 二进制协议(https://github.com/davyxu/goobjfmt)

       内存流直接序列化, 适用于服务器内网传输

    可以通过codec包自行添加新的编码格式

* 支持混合编码收发

    无需改动代码，只需调整消息注册方式，即可达成运行期同时收发不同编码的封包

    - 与Unity3D+Lua使用sproto通信

    - 与其他语言编写的服务器使用protobuf

    - 与web服务器使用json通信

    优点：

    - 享受不同通信协议格式的优点，避免缺点。如Protobuf适合服务器间的传输，封包小，适应面广。

    - 私有的二进制协议方便加密和客户端处理，低GC和实现复杂度。

## 队列实现各种单线程/多线程，异步/同步业务处理模型

  使用队列消息处理结构，方便实现以下效果：

* 单线程异步逻辑，适用于MMORPG复杂交互，免加锁处理共享数据。

* 多线程同步逻辑，适用于机器人逻辑，每个机器人使用独立的goroutine收发处理，机器人间互不干扰。

* 多线程并发逻辑，适用于网关，消息转发，HTTP服务器，每条消息处理都会在完全并发下。

## 远程过程调用（RPC）

* 支持同步RPC，适用于后台服务器向其他服务器请求数据后再顺处理事务。

* 支持异步RPC，适用于单线程服务器逻辑。

## 消息日志
* 可以方便的通过日志查看收发消息的每一个字段消息

![cellnetlogo](doc/chatlogger.png)

# 第三方库依赖

* github.com/davyxu/golog

* github.com/davyxu/goobjfmt

* github.com/golang/protobuf (选择pb编码时依赖)

# 获取+编译

```
	go get -v github.com/davyxu/cellnet

```


# 样例
```golang

const peerAddress = "127.0.0.1:17701"

// 服务器逻辑
func server() {

    // 创建服务器的事件队列，所有的消息，事件都会被投入这个队列处理
	queue := cellnet.NewEventQueue()

    // 创建一个服务器的接受器(Acceptor)，接受客户端的连接
	peerIns := peer.NewGenericPeer("tcp.Acceptor", "server", peerAddress, queue)

    // 将接受器Peer与tcp.ltv的处理器绑定，并设置事件处理回调
    // tcp.ltv处理器负责处理消息收发，使用私有的封包格式以及日志，RPC等处理
	proc.BindProcessorHandler(peerIns, "tcp.ltv", func(ev cellnet.Event) {

        // 处理Peer收到的各种事件
		switch msg := ev.Message().(type) {
		case *cellnet.SessionAccepted: // 接受一个连接
			fmt.Println("server accepted")
		case *TestEchoACK: // 收到连接发送的消息

			fmt.Printf("server recv %+v\n", msg)

            // 发送回应消息
			ev.Session().Send(&TestEchoACK{
				Msg:   msg.Msg,
				Value: msg.Value,
			})

		case *cellnet.SessionClosed: // 会话连接断开
			fmt.Println("session closed: ", ev.Session().ID())
		}

	})

    // 启动Peer，服务器开始侦听
	peerIns.Start()

    // 开启事件队列，开始处理事件，此函数不阻塞
	queue.StartLoop()
}

// 模拟客户端逻辑
func client() {

    // 例子专用的完成标记
	done := make(chan struct{})

    // 创建客户端的事件处理队列
	queue := cellnet.NewEventQueue()

    // 创建客户端的连接器
	peerIns := peer.NewGenericPeer("tcp.Connector", "client", peerAddress, queue)

    // 将客户端连接器Peer与tcp.ltv处理器绑定，并设置接收事件回调
	proc.BindProcessorHandler(peerIns, "tcp.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected: // 已经连接上
			fmt.Println("client connected")
			ev.Session().Send(&TestEchoACK{
				Msg:   "hello",
				Value: 1234,
			})
		case *TestEchoACK: //收到服务器发送的消息

			fmt.Printf("client recv %+v\n", msg)

			// 完成操作
			done <- struct{}{}

		case *cellnet.SessionClosed:
			fmt.Println("client closed")
		}
	})

    // 开启客户端Peer
	peerIns.Start()

    // 开启客户端队列处理
	queue.StartLoop()

	// 等待客户端收到消息
	<-done
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

    chat            聊天

    echo            回音服务器

    fileserver      使用cellnet内建HTTP服务器支持文件服务

msglog              消息日志处理

peer                各种协议的端实现，以及端注册入口及复用组件

    http            HTTP协议处理流程及端封装

    tcp             TCP协议处理流程及端封装

    udp             UDP协议处理流程及端封装

proc                各种处理器实现，以及处理器注册入口

    http            HTTP消息处理及文件服务实现

    tcp             在TCP peer上构建的tcp处理器集合

    udp             在UDP peer上构建的udp处理器集合

relay               接力消息封装

rpc                 远程过程调用支持

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

# 概念及说明

通过下面的介绍，可以深入了解cellnet的组件及构成原理。阅读并理解这些文字，有助于完全上手cellnet的实际使用。

## 队列

队列在cellnet中使用cellnet.Queue接口, 底层由带缓冲的channel实现

在cellnet中, 队列根据实际逻辑需要定制数量. 但一般情况下, 推荐使用一个队列（单线程）处理逻辑。

多线程处理逻辑并不会让逻辑处理更快，过多的同步锁反而会让并发竞态问题变的很严重，导致性能下降严重，同时逻辑编写难度上升。

出现耗时任务时，应该使用生产者和消费者模型，生产者将任务通过channel投放给另外一个goroutine中的消费者处理。

### 创建和开启队列

队列使用NewEventQueue创建，使用.StartLoop()开启队列事件处理循环，所有投递到队列中的函数回调会在队列自由的goroutine中被调用，逻辑在此时被处理

一般在main goroutine中调用queue.Wait阻塞等待队列结束。

```golang
    queue := cellnet.NewEventQueue()

    // 启动队列
    queue.StartLoop()

    // 这里添加队列使用代码

    // 等待队列结束, 调用queue.StopLoop(0)将退出阻塞
    queue.Wait()
```


### 往队列中投递回调
队列中的每一个元素为回调，使用queue的Post方法将回调投递到队列中，回调在Post调用时不会马上被调用。

```golang
    queue.Post(func() {
		fmt.Println("hello")
	})

```

在cellnet正常使用中，Post方法会被封装到内部被调用。正常情况下，逻辑处理无需主动调用queue.Post方法。


## 侦听和接受连接

cellnet使用Acceptor接收多个连接，Acceptor是一种Peer（端），连接到Acceptor的Peer叫做Connector。

一个Peer拥有很多属性（名称，地址，队列），peer.NewGenericPeer函数封装了属性的设置过程。

peer.NewGenericPeer创建好的Peer不会产生任何socket操作，对于Acceptor来说，调用Acceptor的Start方法后，才会真正开始socket的侦听

使用如下代码创建一个接受器(Acceptor)：

```golang
    queue := cellnet.NewEventQueue()

    // NewGenericPeer参数依次是: peer类型, peer名称(日志中方便查看), 侦听地址，事件队列
    peerIns := peer.NewGenericPeer("tcp.Acceptor", "server", "127.0.0.1:8801", queue)

    peerIns.Start()
```


## 创建并发起连接

Connector也是一种Peer，与Acceptor很很多类似的地方，因此创建过程也是类似的。

使用如下代码创建一个连接器(Connector)：

```golang
    queue := cellnet.NewEventQueue()

    peerIns := peer.NewGenericPeer("tcp.Connector", "client", "127.0.0.1:8801", queue)

    peerIns.Start("127.0.0.1:8801")
```

### 自动重连机制
使用golang接口查询特性，可以在peerIns(Peer或GenericPeer接口类型)中查询TCPConnector接口。

该接口可以使用TCPConnector的进一步功能，例如：自动重连。

在服务器连接中，自动重连特性是非常方便的，在连接不成功或者断开时，自动重连会等待一定时间再次发起连接，使用SetReconnectDuration方法可以设置。

```golang
    // 在peerIns接口中查询TCPConnector接口，设置连接超时2秒后自动重连
    peerIns.(cellnet.TCPConnector).SetReconnectDuration(2*time.Second)
```

无需自动重连时，可以使用SetReconnectDuration(0)


## 接收消息

cellnet使用Processor处理消息的收发过程。

使用proc.BindProcessorHandler函数，将一个Peer绑定到某个Processor上，且设置用户消息处理回调。

下面代码尝试将peerIns的Peer，绑定"tcp.ltv"处理器，回调函数为func(ev cellnet.Event) { ... }
    
```golang
proc.BindProcessorHandler(peerIns, "tcp.ltv", func(ev cellnet.Event) {

	switch msg := ev.Message().(type) {
	// 有新的连接连到8801端口
	case *cellnet.SessionAccepted:
		log.Debugln("server accepted")
	// 有连接从8801端口断开
	case *cellnet.SessionClosed:
		log.Debugln("session closed: ", ev.Session().ID())
	// 收到某个连接的ChatREQ消息
	case *proto.ChatREQ:

		// 准备回应的消息
		ack := proto.ChatACK{
			Content: msg.Content,       // 聊天内容
			Id:      ev.Session().ID(), // 使用会话ID作为发送内容的ID
		}

		// 在Peer上查询SessionAccessor接口，并遍历Peer上的所有连接，并发送回应消息（即广播消息）
		p.(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {

			ses.Send(&ack)

			return true
		})

	}

})
```

## 接收系统事件

cellnet将系统事件与收到的消息合并使用用户回调通知逻辑。

下面列出常用的系统事件, 在sysmsg.go文件中定义。

适用Peer类型 | 事件类型 | 事件对应消息
---|---|---
tcp.Connector | 连接成功 | cellnet.SessionConnected
tcp.Connector | 连接错误 | cellnet.SessionConnectError
tcp.Acceptor | 接受新连接 | cellnet.SessionAccepted
tcp.Acceptor/tcpConnector | 会话关闭 | cellnet.SessionClosed


## 发送消息

发送消息往往发生在收到消息或系统事件时，例如：连接上服务器时，发送消息；收到客户端的消息时发送消息。


TCPConnector某些时候需要主动发送消息时，可以这样写
```golang
peerIns.(cellnet.TCPConnector).Session().Send( &YourMsg{ ... } )
```

- 不要缓存Event

cellnet.Event是消息处理的上下文, 可能在底层存在内存池及各种重用行为, 因此不要缓存Event


## 定制自己的Codec
cellnet内建提供基本的编码格式，如果有新的编码需要增加时，可以将这些编码注册到cellnet中。

定制一个自己的Codec，可以直接参考codec/json包下的例子即可。

## 定制自己的Peer
cellnet内建提供的tcp/udp/http能满足90%的Peer需求，但在有些情况下，仍然需要定制新的Peer。

**定制Peer的根本目的：让事件收发处理使用统一的接口和流程**

例如：

- cellnet v4版本暂时没有支持websocket的Peer，可以选定一个第三方库，封装定制为自己的Peer，让Websocket的消息收发与tcp协议一模一样。

- Redis或MySQL连接器可以定制为特殊的Peer，通过统一的Peer Start配合地址就可以方便的发起连接

## cellnet内建Peer类型

Peer类型 | 对应接口 | 功能
---|---|---
tcp.Connector | TCPConnector | tcp发起连接，自动重连
tcp.Acceptor | TCPAcceptor | tcp接受连接，优雅重启
http.Connector | HTTPConnector | http发起请求和接收解码回应
http.Acceptor | HTTPAcceptor | http文件服务，消息收发
udp.Connector | 没有特殊接口 | udp发起连接，无握手
udp.Acceptor | 没有特殊接口 | udp连接管理


# 定制自己的Processor
cellnet提供的Processor能满足基本的通信封包格式的处理，但在特殊封包定制需求时，仍然需要编写自己的Processor

**定制Processor的根本目的：复用业务逻辑层与协议通信层之间的逻辑**

例如：
- 编写服务器逻辑时，已有Unity3D客户端网络库及封包格式，需要服务器与客户端通信，封包格式为私有协议。

- 在连接建立后，需要有握手过程（加密秘钥交换，对时等），此时将这个过程封装到Processor中，可以方便Peer的快速复用。

- 服务器互相连接时，需要标识连接方的服务类型。

- KCP协议需要建立在UDP协议层，UDP收发过程已由Peer部分完成，KCP协议解析只需要放在Processor即可。

- RPC（远程过程调用）是建立在tcp协议之上的调用封装，也可以使用Processor来完成封装过程。

- 消息收发量统计。

- 使用nsq、mysql的Peer的基础上，还需要对用户消息进行二次封装，以实现扩展消息。

## 不需要使用Processor扩展的情况
 
- 编码层更换
   
   从二进制编码更换为ProtocolBuffer编码，从JSON更换为二进制编码等。这种情况下直接使用codec编码包即可。
  
- 具体的业务逻辑


## cellnet内建的Processor列表
Processor类型 | 功能
---|---
tcp.ltv | TCP协议，Length-Type-Value封包格式，带RPC,Relay功能
udp.ltv | UDP协议，Length-Type-Value封包格式
http | 基本HTTP处理



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

# 贡献者

按贡献时间排序，越靠前表示越新的贡献

superikw(https://github.com/superikw), 在v3中测试出一个websocket接口并发发送问题，wss支持，修复会话管理。

bruce.hu(https://github.com/hxdhero), 在v3中测试出一个竞态冲突的bug

M4tou(https://github.com/mutousay), 在v3中协助解决RPC异步超时回调处理

chuan.li(https://github.com/blade-226), 在v3中提供一个没有在io线程编码的bug

Chris Lonng(https://github.com/lonnng), 在v3中提供一个最大封包约束造成服务器间连接断开的bug

IronsDu(https://github.com/IronsDu), 在v2中大幅度性能优化

viwii(viwii@sina.cn), 在v2中，提供一个可能造成死锁的bug


# 备注

感觉不错请star, 谢谢!

开源讨论群: 527430600

知乎: http://www.zhihu.com/people/sunicdavy

提交bug及特性: https://github.com/davyxu/cellnet/issues
