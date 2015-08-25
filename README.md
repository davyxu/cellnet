# CellNet
A Golang game server framework based on actor model

# Target

Erlang like API style

More easy when build game servers

Easy to handle and management

High scalability

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

		done <- true
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
blog: http://www.cppblog.com/sunicdavy

zhihu follow me: http://www.zhihu.com/people/xu-bo-62-87

qq group: 309800774 加群请说明github

mail: sunicdavy@qq.com
