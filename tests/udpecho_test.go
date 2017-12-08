package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/comm/udppeer"
	"github.com/davyxu/cellnet/comm/udppkt"
	"github.com/davyxu/cellnet/tests/proto"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const udpEchoAddress = "127.0.0.1:7701"

var udpEchoSignal *util.SignalTester

var udpEchoAcceptor cellnet.Peer

func StartUDPEchoServer() {

	udpEchoAcceptor = cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "tv.udp.Acceptor",
		PeerAddress: udpEchoAddress,
		PeerName:    "server",
		Event: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *udppkt.SessionAccepted:
					fmt.Println("server accepted")
				case *proto.TestEchoACK:

					fmt.Printf("server recv %+v\n", msg)

					ev.Ses.Send(&proto.TestEchoACK{
						Msg:   msg.Msg,
						Value: msg.Value,
					})

				case *udppkt.SessionClosed:
					fmt.Println("server error: ")
				}
			}

			return nil
		},
	}).Start()

}

func StartUDPEchoClient() {
	queue := cellnet.NewEventQueue()

	cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "tv.udp.Connector",
		Queue:       queue,
		PeerAddress: udpEchoAddress,
		PeerName:    "client",
		Event: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *udppkt.SessionConnected:
					fmt.Println("client connected")
					ev.Ses.Send(&proto.TestEchoACK{
						Msg:   "hello",
						Value: 1234,
					})
				case *proto.TestEchoACK:

					fmt.Printf("client recv %+v\n", msg)

					udpEchoSignal.Done(1)

				case *udppkt.SessionClosed:
					fmt.Println("client error: ")
				}
			}

			return nil
		},
	}).Start()

	queue.StartLoop()

	udpEchoSignal.WaitAndExpect("not recv data", 1)
}

func TestUDPEcho(t *testing.T) {

	udpEchoSignal = util.NewSignalTester(t)

	StartUDPEchoServer()

	StartUDPEchoClient()

	udpEchoAcceptor.Stop()
}

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
