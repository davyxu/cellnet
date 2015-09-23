package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/golang/protobuf/proto"
	"log"
	"runtime"
	"strconv"
	"sync"
)

var done = make(chan bool)

// 测试客户端连接数量
const connCount = 10

func runClient() {

	evq := cellnet.NewEvQueue()

	// 同步量
	var endAcc sync.WaitGroup

	socket.RegisterMessage(evq, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

		// 正常收到
		endAcc.Done()
	})

	socket.RegisterMessage(evq, coredef.ConnectedACK{}, func(ses cellnet.Session, content interface{}) {

		id, _ := strconv.Atoi(ses.FromPeer().Name())

		// 连接上发包
		ses.Send(&coredef.TestEchoACK{
			Content: proto.String(fmt.Sprintf("data#%d", id)),
		})

	})

	// 启动N个连接
	for i := 0; i < connCount; i++ {

		endAcc.Add(1)

		socket.NewConnector(evq).Start("127.0.0.1:7235").SetName(fmt.Sprintf("%d", i))

	}

	log.Println("waiting server msg...")

	// 等待完成
	endAcc.Wait()

}

func runServer() {
	evq := cellnet.NewEvQueue()

	acc := socket.NewAcceptor(evq).Start("127.0.0.1:7235").(cellnet.SessionManager)

	socket.RegisterMessage(evq, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		if acc.Get(ses.ID()) != ses {
			panic("1: session not exist in SessionManager")
		}

		log.Println("server recv:", msg.String())

		// 发包后关闭
		ses.Send(&coredef.TestEchoACK{
			Content: proto.String(msg.GetContent()),
		})

		if acc.Get(ses.ID()) != ses {
			panic("2: session not exist in SessionManager")
		}

		ses.Close()

		if acc.Get(ses.ID()) != ses {
			panic("3: session not exist in SessionManager")
		}

	})

	done <- true

}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	mode := flag.String("mode", "", "specify the mode of this test")

	flag.Parse()

	if mode != nil && *mode == "client" {
		runClient()
	} else {
		runServer()
	}

}
