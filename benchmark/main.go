package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/json"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
	"log"
	"os"
	"reflect"
	"runtime/pprof"
	"time"
)

func server() {
	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Acceptor", "server", "127.0.0.1:7701", queue)

	dispatcher := proc.NewMessageDispatcher(p, "tcp.ltv")

	dispatcher.RegisterMessage("main.TestEchoACK", func(ev cellnet.Event) {

		msg := ev.Message().(*TestEchoACK)

		ev.Session().Send(&TestEchoACK{
			Msg:   msg.Msg,
			Value: msg.Value,
		})
	})

	p.Start()

	queue.StartLoop()
}

func client() {

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "client", "127.0.0.1:7701", queue)

	proc.BindProcessorHandler(p, "tcp.ltv", nil)

	p.Start()

	queue.StartLoop()

	rv := proc.NewSyncReceiver(p)

	rv.Recv(func(ev cellnet.Event) {
		msg := ev.Message().(*cellnet.SessionConnected)
		msg = msg

		ev.Session().Send(&TestEchoACK{
			Msg:   "hello",
			Value: 1234,
		})

	})

	begin := time.Now()

	for time.Now().Sub(begin) < 10*time.Second {

		rv.Recv(func(ev cellnet.Event) {

			ev.Session().Send(&TestEchoACK{
				Msg:   "hello",
				Value: 1234,
			})

		})
	}

}

var profile = flag.String("profile", "", "write cpu profile to file")

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    int(util.StringHash("main.TestEchoACK")),
	})
}

// go build -o bench.exe main.go
// ./bench.exe -profile=mem.pprof
// go tool pprof -alloc_space -top bench.exe mem.pprof
func main() {

	flag.Parse()

	f, err := os.Create(*profile)
	if err != nil {
		log.Fatal(err)
	}

	golog.SetLevelByString("*", "info")

	server()

	client()

	pprof.WriteHeapProfile(f)
	f.Close()
}
