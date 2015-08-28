# CellNet
A Golang game server framework based on actor model

# Feature

本机,跨进程,跨机器通信均使用统一的Actor模型

为游戏服务器优化, 注重开发效率及运行效率

# Roadmap
RPC支持

服务器框架例子

服务器可视化部署工具

# Dependencies
github.com/golang/protobuf/proto

github.com/BurntSushi/toml

# Example
=================================
## Hello world
```go


cid := cellnet.Spawn(func(_ cellnet.CellID, cl interface{}) {

	switch v := cl.(type) {
	case string:
		log.Println(v)
	}

})

cellnet.Send(cid, "hello world ")


```

## Client & server with message dispatcher
```go
func server() {

	disp := dispatcher.NewPacketDispatcher()

	dispatcher.RegisterMessage(disp, coredef.TestEchoACK{}, func(ses cellnet.CellID, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("server recv:", msg.String())

		cellnet.Send(ses, &coredef.TestEchoACK{
			Content: proto.String("world"),
		})
	})

	ltvsocket.SpawnAcceptor("127.0.0.1:8001", dispatcher.PeerHandler(disp))
}

func client() {

	disp := dispatcher.NewPacketDispatcher()

	dispatcher.RegisterMessage(disp, coredef.TestEchoACK{}, func(ses cellnet.CellID, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

	})

	dispatcher.RegisterMessage(disp, coredef.ConnectedACK{}, func(ses cellnet.CellID, content interface{}) {
		cellnet.Send(ses, &coredef.TestEchoACK{
			Content: proto.String("hello"),
		})
	})

	ltvsocket.SpawnConnector("127.0.0.1:8001", dispatcher.PeerHandler(disp))

}

```

# Contact 
博客: http://www.cppblog.com/sunicdavy

知乎: http://www.zhihu.com/people/xu-bo-62-87

技术讨论组: 309800774 加群请说明cellnet

邮箱: sunicdavy@qq.com
