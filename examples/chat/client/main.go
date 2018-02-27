package main

import (
	"bufio"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/examples/chat/proto"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
	"os"
	"strings"

	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

var log = golog.New("client")

func ReadConsole(callback func(string)) {

	for {
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}

		text = strings.TrimRight(text, "\r\n ")

		text = strings.TrimLeft(text, " ")

		callback(text)

	}

}

func main() {

	queue := cellnet.NewEventQueue()

	p := peer.NewPeer("tcp.Connector")
	pset := p.(cellnet.PropertySet)
	pset.SetProperty("Address", "127.0.0.1:8801")
	pset.SetProperty("Name", "client")
	pset.SetProperty("Queue", queue)

	proc.BindProcessor(p, "tcp.ltv", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:
			log.Debugln("client connected")
		case *cellnet.SessionClosed:
			log.Debugln("client error")
		case *proto.ChatACK:
			log.Infof("sid%d say: %s", msg.Id, msg.Content)
		}
	})

	p.Start()

	queue.StartLoop()

	ReadConsole(func(str string) {

		p.(interface {
			Session() cellnet.Session
		}).Session().Send(&proto.ChatREQ{
			Content: str,
		})

	})

}
