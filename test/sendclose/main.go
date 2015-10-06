package main

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
)

var done = make(chan int)

// 测试客户端连接数量
const connCount = 10

// 多连接收封包后被服务器关闭, 确保收到封包
func multiConn() {

	pipe := cellnet.NewEventPipe()

	// 同步量
	var endAcc sync.WaitGroup

	// 启动N个连接
	for i := 0; i < connCount; i++ {

		endAcc.Add(1)

		p := socket.NewConnector(pipe).Start("127.0.0.1:7235")

		p.SetName(fmt.Sprintf("%d", i))

		socket.RegisterSessionMessage(p, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
			msg := content.(*coredef.TestEchoACK)

			log.Println("client recv:", msg.String())

			// 正常收到
			endAcc.Done()
		})

		socket.RegisterSessionMessage(p, coredef.SessionConnected{}, func(ses cellnet.Session, content interface{}) {

			id, _ := strconv.Atoi(ses.FromPeer().Name())

			// 连接上发包
			ses.Send(&coredef.TestEchoACK{
				Content: proto.String(fmt.Sprintf("data#%d", id)),
			})

		})

	}

	pipe.Start()

	// 等待完成
	endAcc.Wait()

	fmt.Println("multi connection close test done!")

}

// 客户端连接上后, 主动断开连接, 确保连接正常关闭
func connClose() {

	pipe := cellnet.NewEventPipe()

	p := socket.NewConnector(pipe).Start("127.0.0.1:7235")

	socket.RegisterSessionMessage(p, coredef.SessionConnected{}, func(ses cellnet.Session, content interface{}) {

		// 连接上发包,告诉服务器不要断开
		ses.Send(&coredef.TestEchoACK{
			Content: proto.String("noclose"),
		})

	})

	socket.RegisterSessionMessage(p, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())
		done <- 1

		// 客户端主动断开
		ses.Close()

	})

	socket.RegisterSessionMessage(p, coredef.SessionClosed{}, func(ses cellnet.Session, content interface{}) {

		log.Println("close ok!")
		// 正常断开
		done <- 2

	})

	pipe.Start()

	// 收到回包
	if <-done != 1 {
		log.Panicln("test failed, not recv msg")
	}

	// 断开正常
	if <-done != 2 {
		log.Panicln("test failed, not close")
	}

	fmt.Println("connected close test done!")

}

func runServer() {
	pipe := cellnet.NewEventPipe()

	p := socket.NewAcceptor(pipe).Start("127.0.0.1:7235")

	// 计数器, 应该按照connCount倍数递增
	var counter int

	socket.RegisterSessionMessage(p, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		counter++
		log.Printf("No. %d: server recv: %v", counter, msg.String())

		// 发包后关闭
		ses.Send(&coredef.TestEchoACK{
			Content: proto.String(msg.GetContent()),
		})

		if msg.GetContent() != "noclose" {
			ses.Close()
		}

	})

	pipe.Start()

	done <- 0

}

// 运行服务器: sendclose s
// 运行客户端: sendclose
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) > 1 && os.Args[1] == "s" {
		runServer()
	} else {
		multiConn()
		connClose()
	}

}
