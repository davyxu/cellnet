package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	_ "github.com/davyxu/cellnet/comm/udppeer"
	_ "github.com/davyxu/cellnet/comm/udppkt"
	"github.com/davyxu/cellnet/tests/proto"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const udpSeqAddress = "127.0.0.1:7901"

var udpSeqSignal *util.SignalTester

var udpSeqAcceptor cellnet.Peer

func StartUDPSeqServer() {

	udpSeqAcceptor = cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "tv.udp.Acceptor",
		PeerAddress: udpSeqAddress,
		PeerName:    "server",
		Event: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *comm.SessionAccepted:
					fmt.Println("server accepted")
				case *proto.TestEchoACK:

					//fmt.Printf("server recv %+v\n", msg)

					ev.Ses.Send(&proto.TestEchoACK{
						Msg:   msg.Msg,
						Value: msg.Value,
					})

				case *comm.SessionClosed:
					fmt.Println("server error: ")
				}
			}

			return nil
		},
	}).Start()

}

func StartUDPSeqClient() {

	var counter int32

	cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "tv.udp.Connector",
		PeerAddress: udpSeqAddress,
		PeerName:    "client",
		Event: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *comm.SessionConnected:
					fmt.Println("client connected")
					ev.Ses.Send(&proto.TestEchoACK{
						Value: counter,
					})

				case *proto.TestEchoACK:

					if msg.Value != counter {
						fmt.Println("seq not match")
						udpSeqSignal.FailNow()
					}

					counter++
					ev.Ses.Send(&proto.TestEchoACK{
						Value: counter,
					})

				case *comm.SessionClosed:
					fmt.Println("client error: ")
				}
			}

			return nil
		},
	}).Start()

}

func TestUDPSeq(t *testing.T) {

	udpSeqSignal = util.NewSignalTester(t)

	StartUDPSeqServer()

	StartUDPSeqClient()

	queue := cellnet.NewEventQueue()

	queue.StartLoop()
	queue.Wait()

	udpSeqAcceptor.Stop()
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
