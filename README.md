# cellnet

[![Build Status](https://api.travis-ci.org/davyxu/cellnet.svg?branch=v3)](https://travis-ci.org/davyxu/cellnet)

简单,方便,高效的跨平台服务器网络库


# 特性

## 队列及IO
  
* 支持多个队列, 实现单线程/多线程收发处理消息

* 多线程处理io

* 发送时自动合并封包(性能效果决定于实际请求和发送比例)

## 数据协议

* 编码支持:
    - Google Protobuf (https://github.com/google/protobuf)
    - sproto (https://github.com/cloudwu/sproto)
    - json
    - 二进制协议(https://github.com/davyxu/goobjfmt)

* 支持混合编码收发

* 传输协议支持:

   - tcp(基于Type-Length-Value私有协议)

   - WebSocket


## 基于handler无状态处理链

* 自定义, 组装收发流程

* 支持专有日志调试

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

# websocket可选支持

* github.com/gorilla/websocket

# 获取+编译

```
	go get -u -v github.com/davyxu/cellnet
```

# 性能测试

命令行: go test -v github.com/davyxu/cellnet/benchmark/io

平台: Windows 7 x64/CentOS 6.5 x64

测试用例: localhost 1000连接 同时对服务器进行实时PingPong测试

配置1: i7 6700 3.4GHz 8核

IOPS: 12.5w

配置2: i5 4590 3.3GHz 4核

IOPS: 10.1w

# 例子
## Echo
```go


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

	cellnet.RegisterMessage(p, "coredef.SessionConnectFailed", func(ev *cellnet.SessionEvent) {

		msg := ev.Msg.(*coredef.SessionConnectFailed)

		log.Debugln(msg.Reason)

	})

	queue.StartLoop()
}

```

# 文件夹功能

```
benchmark\		    性能测试用例

proto\			    cellnet内部的proto

    binary\         内部系统消息,rpc消息协议

    pb\             使用pb例子的消息

    sproto\         使用sproto例子的消息

protoc-gen-msg\     protobuf的protoc插件, 消息id生成, 使用pb编码时使用

rpc\			    异步远程过程调用封装

socket\			    套接字,连接管理等封装

example\			测试用例/例子

    classicrecv\    传统的固定消息处理函数例子
	
   	echo_pb\	    基于protobuf和json混合编码的pingpong测试，

   	echo_sproto\	基于sproto编码的pingpong测试，

   	echo_websocket\	基于websocket协议，

   	gracefulexit\	平滑退出

	rpc\		    异步远程过程调用

	sendclose\		发送消息后保证消息送达后再断开连接
	
	timer\		    异步计时器
	

timer\			计时器接口

util\			工具库

```

# FAQ

* 这个代码的入口在哪里? 怎么编译为exe?

    本代码是一个网络库, 需要根据需求, 整合逻辑

    只需要将sample里echo系代码复制到你的main中编译即可运行

* 支持WebSocket么?

    支持!

    本网络库的Websocket基于第三方整合, 包格式基于文本: 包名\n+json内容

    tcp私有协议到Websocket的转换, 只需要更换包名即可

* 混合编码有何用途?

    在与多种语言写成的服务器进行通信时, 可以使用不同的编码,
    最终在逻辑层都是统一的结构能让逻辑编写更加方便, 无需关注底层处理细节

* 内建支持的二进制协议能与其他语言写成的网络库互通么?

    完全支持, 但内建二进制协议支持更适合网关与后台服务器.
    不建议与客户端通信中使用, 二进制协议不会忽略使用默认值的字段

* 我能通过Handler处理链进行怎样的扩展?

    封包需要加密, 统计, 预处理时, 可以使用Handler. 每个Handler建议无状态,
    需要存储的数据, 可以通过Event中的Tag进行扩展

* 如何查看Handler处理流程?
    在程序启动时, 调用如下代码
```
    cellnet.EnableHandlerLog = true
```

   可在日志中看到如下日志格式

```
    [DEBUG] cellnet 2000/00/00 01:02:03 9 Event_Connected [svc->agent] <DecodePacketHandler> SesID: 1 MsgID: 3551021301(coredef.SessionConnected) {} Tag: <nil> TransmitTag: <nil> Raw: (0)[]
```

    9 表示一个Event处理序号, 同一序号表示1个处理流程, 例如1个接收/发送流程

    Event_Connected 表示事件名

    [svc->agent] 表示peer的名称

    <DecodePacketHandler> 表示Handler的名称, 通过反射取得

    SesID 表示 会话ID, 由SessionManager分配

    MsgID 表示消息号, 后面括号中是对应的消息名, 如果未在系统中注册, 显示为空, 后续是消息内容

    TransmitTag, Tag 附属上下文内容

    Raw, 表示消息的原始二进制信息


* 所有的例子都是单线程的, 能编写多线程的逻辑么?

    完全可以, cellnet并没有全局的队列, 只需在Acceptor和Connector创建时,
    传入不同的队列, socket收到的消息就会被放到这个队列中
    传入空队列时, 使用并发方式(io线程)调用处理回调

* 消息日志为什么与处理函数日志顺序不统一?

    由于消息日志反应的是收到消息的日志, 因此必须放置在io线程中处理. 而单线程逻辑与io线程分别在不同的线程. 日志顺序错位是正常的
    如果需要顺序日志: 可以在进程启动时, 调用runtime.GOMAXPROCS(1), 将go的线程调度默认为1CPU

* cellnet有网关和db支持么?

    cellnet专注于服务器底层.你可以根据自己需要编写网关及db支持



# 版本历史
2017.1  v2版本 [详细请查看](https://github.com/davyxu/cellnet/blob/master/CHANGES.md)

2015.8	v1版本


# 贡献者

viwii(viwii@sina.cn), 提供一个可能造成死锁的bug

IronsDu(duzhongwei@qq.com), 大幅度性能优化

Chris Lonng(chris@lonng.org), 提供一个最大封包约束造成服务器间连接断开的bug

# 备注

感觉不错请star, 谢谢!

博客: http://www.cppblog.com/sunicdavy

知乎: http://www.zhihu.com/people/sunicdavy

提交bug及特性: https://github.com/davyxu/cellnet/issues
