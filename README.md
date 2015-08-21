# CellNet
A Golang game server framework based on actor model

# Target

Erlang like API style

More easy when build game servers

Ez to handle and management

High scalability

# Dependencies
github.com/golang/protobuf/proto


# Example
=================================
```go

cid := cellnet.Spawn(func(mailbox chan interface{}) {
	for {

		switch v := (<-mailbox).(type) {
		case string:
			log.Println(v)
		}
	}

})

cellnet.Send(cid, "hello world ")


```

# Contact 
blog: http://www.cppblog.com/sunicdavy

zhihu follow me: http://www.zhihu.com/people/xu-bo-62-87

qq group: 309800774 加群请说明github

mail: sunicdavy@qq.com
