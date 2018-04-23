# 定制Peer
cellnet内建提供的tcp/udp/http能满足90%的Peer需求，但在有些情况下，仍然需要定制新的Peer。

**定制Peer的根本目的：让事件收发处理使用统一的接口和流程**

例如：

- cellnet v4版本暂时没有支持websocket的Peer，可以选定一个第三方库，封装定制为自己的Peer，让Websocket的消息收发与tcp协议一模一样。

- Redis或MySQL连接器可以定制为特殊的Peer，通过统一的Peer Start配合地址就可以方便的发起连接