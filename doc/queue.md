# 队列

队列在cellnet中使用cellnet.Queue接口, 底层由带缓冲的channel实现


## 创建和开启队列

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


## 往队列中投递回调
队列中的每一个元素为回调，使用queue的Post方法将回调投递到队列中，回调在Post调用时不会马上被调用。

```golang
    queue.Post(func() {
		fmt.Println("hello")
	})

```

在cellnet正常使用中，Post方法会被封装到内部被调用。正常情况下，逻辑处理无需主动调用queue.Post方法。

## 多线程和单线程
cellnet中，一个队列可以理论上对应一个线程。默认所有例子都是单队列单线程，这种处理方法并不慢。

在cellnet中, 队列根据实际逻辑需要定制数量. 但一般情况下, 推荐使用一个队列（单线程）处理逻辑。

多线程处理逻辑并不会让逻辑处理更快，过多的同步锁反而会让并发竞态问题变的很严重，导致性能下降严重，同时逻辑编写难度上升。

出现耗时任务时，应该使用生产者和消费者模型，生产者将任务通过channel投放给另外一个goroutine中的消费者处理。

需要多线程并发处理时，请在所有peer需要传入队列的地方设置为nil。消息将在IO线程中被派发并推给逻辑层处理。