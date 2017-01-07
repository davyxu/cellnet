# V2版本

## 变化及修改

- 去掉V1中的EventPipe

- V1中的EventQueue被拆成EventDispatcher及新的EventQueue

- 新的EventQueue实现了EventPipe的一部分功能

- 调整EventQueue的Post命名及DelayPost的参数

- 去掉PeerEvent支持

- socket.RegisterEventMessage改为socket.RegisterMessage

- 例子/测试用例使用sample文件夹命名

## V2版本特性

- 实现单线程逻辑时, 全局只有1个EventQueue. 而不是一个Peer一个Queue

- EventDispatcher处理回调

- 处理DB, Timer等不依赖Dispatcher(Peer)逻辑时, 在Post时, Dispatcher可以指定nil, 通过data的函数得到异步返回