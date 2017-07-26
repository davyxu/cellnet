package gracefulexit

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/golog"
	"sync"
	"testing"
	"time"
)

var log *golog.Logger = golog.New("test")

const clientConnectionCount = 3

func TestCreateDestroyAcceptor(t *testing.T) {
	queue := cellnet.NewEventQueue()

	p := socket.NewAcceptor(queue).Start("127.0.0.1:7701")
	p.SetName("server")

	var allAccepted sync.WaitGroup
	cellnet.RegisterMessage(p, "coredef.SessionAccepted", func(ev *cellnet.Event) {

		allAccepted.Done()
	})

	queue.StartLoop()

	log.Debugln("Start connecting...")
	allAccepted.Add(clientConnectionCount)
	runMultiConnection()

	log.Debugln("Wait all accept...")
	allAccepted.Wait()

	log.Debugln("Close acceptor...")
	p.Stop()

	// 确认所有连接已经断开
	time.Sleep(time.Second)

	log.Debugln("Session count:", p.SessionCount())

	p.Start(p.Address())
	log.Debugln("Start connecting...")
	allAccepted.Add(clientConnectionCount)
	runMultiConnection()

	log.Debugln("Wait all accept...")
	allAccepted.Wait()

	log.Debugln("All done")
}

func runMultiConnection() {

	for i := 0; i < clientConnectionCount; i++ {

		p := socket.NewConnector(nil).Start("127.0.0.1:7701")
		p.SetName("client.MultiConn")
	}

}
