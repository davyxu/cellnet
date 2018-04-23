# 收发及处理消息

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
## cellnet内建的Processor列表
Processor类型 | 功能
---|---
tcp.ltv | TCP协议，Length-Type-Value封包格式，带RPC,Relay功能
udp.ltv | UDP协议，Length-Type-Value封包格式
http | 基本HTTP处理


## 接收系统事件

cellnet将系统事件使用消息派发给用户，这种消息并不是由Socket接收，而是由cellnet内部生成的。

例如：当TCP socket连接上服务器时，在回调中，将会收到一个*cellnet.SessionConnected的消息

下面列出常用的系统事件, 在sysmsg.go文件中定义。

适用Peer类型 | 事件类型 | 事件对应消息
---|---|---
tcp.Connector | 连接成功 | cellnet.SessionConnected
tcp.Connector | 连接错误 | cellnet.SessionConnectError
tcp.Acceptor | 接受新连接 | cellnet.SessionAccepted
tcp.Acceptor/tcpConnector | 会话关闭 | cellnet.SessionClosed

这样设计的益处：

- 无需为系统事件准备另外的一套处理回调

- 系统事件对应的消息也可以使用Hooker处理或者过滤

## 发送消息

发送消息往往发生在收到消息或系统事件时，例如：连接上服务器时，发送消息；收到客户端的消息时发送消息。


TCPConnector某些时候需要主动发送消息时，可以这样写
```golang
peerIns.(cellnet.TCPConnector).Session().Send( &YourMsg{ ... } )
```

- 不要缓存Event

cellnet.Event是消息处理的上下文, 可能在底层存在内存池及各种重用行为, 因此不要缓存Event