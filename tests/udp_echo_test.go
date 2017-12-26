package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	_ "github.com/davyxu/cellnet/comm/udppeer"
	_ "github.com/davyxu/cellnet/comm/udpproc"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const udpEchoAddress = "127.0.0.1:7901"

var udpEchoSignal *util.SignalTester

var udpEchoAcceptor cellnet.Peer

func StartUDPEchoServer() {

	udpEchoAcceptor = cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Acceptor",
		EventProcessor: "udp.ltv",
		PeerAddress:    udpEchoAddress,
		PeerName:       "server",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *TestEchoACK:

					fmt.Printf("server recv %+v\n", msg)

					ev.Ses.Send(&TestEchoACK{
						Msg:   msg.Msg,
						Value: msg.Value,
					})
				}
			}

			return nil
		},
	}).Start()

}

func StartUDPEchoClient() {

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Connector",
		EventProcessor: "udp.ltv",
		PeerAddress:    udpEchoAddress,
		PeerName:       "client",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *comm.SessionConnected:
					fmt.Println("client connected")
					ev.Ses.Send(&TestEchoACK{
						Msg:   "hello",
						Value: 1234,
					})
				case *TestEchoACK:

					fmt.Printf("client recv %+v\n", msg)

					udpEchoSignal.Done(1)

				case *comm.SessionClosed:
					fmt.Println("client error: ")
				}
			}

			return nil
		},
	}).Start()

	udpEchoSignal.WaitAndExpect("not recv data", 1)
}

func TestUDPEcho(t *testing.T) {

	udpEchoSignal = util.NewSignalTester(t)

	StartUDPEchoServer()

	StartUDPEchoClient()

	udpEchoAcceptor.Stop()
}

//func TestUDPServer(t *testing.T) {
//
//	StartUDPEchoServer()
//
//	queue := cellnet.NewEventQueue()
//
//	queue.StartLoop()
//	queue.Wait()
//}

/*
	_, err = conn.Write([]byte{})

	if err != nil {

		log.Errorf("#write failed(%s) %v", self.NameOrAddress(), err.Error())
		return
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {

		log.Errorf("#read failed(%s) %v", self.NameOrAddress(), err.Error())
		return
	}

	da := binary.BigEndian.Uint32(buff[:n])

	log.Debugln(time.Unix(int64(da), 0).String(), buff[:n])

*/
//const udpAddress = "time.nist.gov:37"
//
//func TestUDPConnector(t *testing.T) {
//
//	queue := cellnet.NewEventQueue()
//
//	cellnet.NewPeer(cellnet.PeerConfig{
//		PeerType:    "udp.Connector",
//		Queue:       queue,
//		PeerAddress: udpAddress,
//		PeerName:    "client",
//	}).Start()
//
//	queue.StartLoop()
//
//	queue.Wait()
//}
