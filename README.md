# cellnet
简单,方便,高效的跨平台服务器网络库


# 特性

## 队列及IO
  
* 支持多个队列, 实现单线程/多线程收发处理消息

* 多线程处理io

* 发送时自动合并封包(性能效果决定于实际请求和发送比例)

## 数据协议

* 封包类型采用Type-Length-Value的私有tcp封包, 自带序列号防御简单的封包复制

* 内建Protobuf, sproto, json, 二进制协议支持

* 支持混合协议收发

## 基于handler处理链, 自定义收发流程

* handler支持日志调试流程

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

	cellnet.RegisterMessage(evd, "gamedef.TestEchoACK", func(ev *cellnet.SessionEvent) {
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

	dh := socket.NewConnector(queue).Start("127.0.0.1:7301")

	cellnet.RegisterMessage(dh, "gamedef.TestEchoACK", func(ev *cellnet.SessionEvent) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.Content)
	})

	cellnet.RegisterMessage(dh, "coredef.SessionConnected", func(ev *cellnet.SessionEvent) {

		log.Debugln("client connected")

		ev.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

	})

	cellnet.RegisterMessage(dh, "coredef.SessionConnectFailed", func(ev *cellnet.SessionEvent) {

		msg := ev.Msg.(*coredef.SessionConnectFailed)

		log.Debugln(msg.Reason)

	})

	queue.StartLoop()
}

```

# 文件夹功能

```
benchmark\		性能测试用例

proto\			cellnet内部的proto
    binary\     内部系统消息,rpc消息协议
    pb\         使用pb例子的消息
    sproto\     使用sproto例子的消息

protoc-gen-msg\ protobuf的protoc插件, 消息id生成

rpc\			异步远程过程调用封装

socket\			套接字,连接管理等封装

example\			测试用例/例子
   
	close\		发送消息后保证消息送达后再断开连接
	
   	echo_pb\	基于protobuf协议的pingpong测试，

   	echo_sproto\	基于sproto协议的pingpong测试，
	
	rpc\		异步远程过程调用
	
	timer\		异步计时器
	

timer\			计时器接口

util\			工具库

```

# FAQ
* 混合协议有何用途?

    在与多种语言写成的服务器进行通信时, 可以使用不同的协议,
    最终在逻辑层都是统一的结构能让逻辑编写更加方便, 无需关注底层处理细节

* 内建支持的二进制协议能与其他语言写成的网络库互通么?

    完全支持, 但内建二进制协议支持更适合网关与后台服务器.
    不建议与客户端通信中使用, 二进制协议不会忽略使用默认值的字段

* 我能通过Handler处理链进行怎样的扩展?

    封包需要加密, 统计, 预处理时, 可以使用Handler. 每个Handler建议无状态,
    需要存储的数据, 可以通过SessionEvent中的Tag进行扩展

* 所有的例子都是单线程的, 能编写多线程的逻辑么?

    完全可以, cellnet并没有全局的队列, 只需在Acceptor和Connector创建时,
    传入不同的队列, socket收到的消息就会被放到这个队列中

* cellnet有网关和db支持么?

    cellnet专注于服务器底层.你可以根据自己需要编写网关及db支持

# 版本历史
2017.1  v2版本 [详细请查看](https://github.com/davyxu/cellnet/blob/master/CHANGES.md)

2015.8	v1版本


# 贡献者

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