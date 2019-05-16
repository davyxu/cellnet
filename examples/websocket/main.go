package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"

	_ "github.com/davyxu/cellnet/codec/json"
	_ "github.com/davyxu/cellnet/peer/gorillaws"
	_ "github.com/davyxu/cellnet/proc/gorillaws"
)

var log = golog.New("websocket_server")

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

// 将消息注册到系统
func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    1234,
	})
}

var (
	flagClient = flag.Bool("client", false, "client mode")
)

const (
	TestAddress = "http://127.0.0.1:18802/echo"
)

func client() {
	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件
	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("gorillaws.Connector", "client", TestAddress, queue)
	p.(cellnet.WSConnector).SetReconnectDuration(time.Second)

	proc.BindProcessorHandler(p, "gorillaws.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {

		case *cellnet.SessionConnected:
			log.Debugln("server connected")

			ev.Session().Send(&TestEchoACK{
				Msg:   "鲍勃",
				Value: 331,
			})
			// 有连接断开
		case *cellnet.SessionClosed:
			log.Debugln("session closed: ", ev.Session().ID())
		case *TestEchoACK:

			log.Debugf("recv: %+v %v", msg, []byte("鲍勃"))

		}
	})

	// 开始侦听
	p.Start()

	// 事件队列开始循环
	queue.StartLoop()

	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	queue.Wait()
}

func server() {
	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	queue := cellnet.NewEventQueue()

	// 侦听在18802端口
	p := peer.NewGenericPeer("gorillaws.Acceptor", "server", TestAddress, queue)

	proc.BindProcessorHandler(p, "gorillaws.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {

		case *cellnet.SessionAccepted:
			log.Debugln("server accepted")
			// 有连接断开
		case *cellnet.SessionClosed:
			log.Debugln("session closed: ", ev.Session().ID())
		case *TestEchoACK:

			log.Debugf("recv: %+v %v", msg, []byte("鲍勃"))

			val, exist := ev.Session().(cellnet.ContextSet).GetContext("request")
			if exist {
				if req, ok := val.(*http.Request); ok {
					raw, _ := json.Marshal(req.Header)
					log.Debugf("origin request header: %s", string(raw))
				}
			}

			ev.Session().Send(&TestEchoACK{
				Msg:   "中文",
				Value: 1234,
			})
		}
	})

	// 开始侦听
	p.Start()

	// 事件队列开始循环
	queue.StartLoop()

	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	queue.Wait()

}

// 默认启动服务器端
// 网页连接服务器： 在浏览器(Chrome)中打开index.html, F12打开调试窗口->Console标签 查看命令行输出
// 	注意：日志中的http://127.0.0.1:18802/echo链接是api地址，不是网页地址，直接打开无法正常工作
// 	注意：如果http代理/VPN在运行时可能会导致无法连接, 请关闭
// 客户端连接服务器：命令行模式中添加-client
func main() {

	flag.Parse()

	if *flagClient {
		client()
	} else {
		server()
	}

}
