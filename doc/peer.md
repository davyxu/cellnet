# 端(Peer)

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

## cellnet内建Peer类型

Peer类型 | 对应接口 | 功能
---|---|---
tcp.Connector | TCPConnector | tcp发起连接，自动重连
tcp.Acceptor | TCPAcceptor | tcp接受连接，优雅重启
http.Connector | HTTPConnector | http发起请求和接收解码回应
http.Acceptor | HTTPAcceptor | http文件服务，消息收发
udp.Connector | UDPConnector | udp发起连接，无握手
udp.Acceptor | 没有特殊接口 | udp连接管理
gorillaws.Acceptor | WSAcceptor | websocket连接管理，加密连接