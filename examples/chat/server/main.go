package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/examples/chat/proto"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"

	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

var log = golog.New("server")

func main() {

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Acceptor", "server", "127.0.0.1:8801", queue)

	proc.BindProcessor(p, "tcp.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *cellnet.SessionAccepted:
			log.Debugln("server accepted")
		case *cellnet.SessionClosed:
			log.Debugln("session closed: ", ev.Session().ID())
		case *proto.ChatREQ:

			ack := proto.ChatACK{
				Content: msg.Content,
				Id:      ev.Session().ID(),
			}

			p.(cellnet.SessionAccessor).VisitSession(func(ses cellnet.Session) bool {

				ses.Send(&ack)

				return true
			})

		}

	})

	p.Start()

	queue.StartLoop()

	queue.Wait()

}
