# FAQ

* 报错panic: Peer type not found怎么办?
    
    这是由于需要的peer没有找到或者没有注册，使用cellnet内建的peer请在main入口这样导入包
```
    import (
        _ "github.com/davyxu/cellnet/peer/tcp"
    )
```

* 报错panic: processor not found怎么办?

    这是由于需要的processor没有找到或者没有注册，使用cellnet内建的processor请在main入口这样导入包
```
    import (
        _ "github.com/davyxu/cellnet/proc/tcp"
    )
```

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

   github.com/davyxu/cellnet/peer/mysql   MySQL支持
   
   github.com/davyxu/cellnet/peer/redis   Redis支持
   
   使用方法请参考tests

* 如何关闭调试消息日志?

   golog.SetLevelByString(".", "info") // 将所有日志的级别提高到info级别，debug低于info级别所以不再显示

   第一个参数支持正则表达式，"."表示所有日志。可以指定日志名关闭

* cellnet能承受多少连接？

   承受连接数量和操作系统和硬件有关系，cellnet本身承载数受操作系统和硬件约束。

* cellnet能做百万请求的服务器么？

   这是架构设计的问题，和cellnet无关。

* 为什么把客户端关掉，没有收到cellnet.SessionClosed事件，内存不降？

   TCP挥手失败不会触发cellnet.SessionClosed，请通过修改peer上的TCPSocketOption接口的SetSocketDeadline，设置读超时避免这个问题。

   游戏服务器请自行实现心跳封包逻辑，以避免攻击者只连接不发包消耗服务器资源。

   TCPSocketOption 接口被TCPAcceptor和TCPConnector实现，因此只要拥有这两种peer都可以直接进行设置，例如：

   // 设置30秒读超时和5秒写超时
   peer.(TCPSocketOption).SetSocketDeadline(time.Second * 30, time.Second * 5)


* cellnet的http能做路由么？能做web服务器么？

   v4版本中添加的http功能是为了方便用通用的方式接收http消息。如果需要专业的http路由，请使用成熟的http服务器，例如gin。

* 为什么发送20k的TCP封包会断开？

   TCP封包请在逻辑层约束到MTU(Maximum Transmission Unit)范围内,一般路由器设置为1500，考虑到包头损耗，一般用户数据大约在1400字节较为安全。

   超过MTU后，在某些路由器将发生封包重传，导致传输性能下降，严重的导致丢包乃至连接断开。

   cellnet底层没有拆分逻辑包的设计，请自行使用Processor扩展。

* 如何获取会话的远程IP?

   util.GetRemoteAddrss获取到地址, util.SpliteAddress拆分出host部分就是ip
