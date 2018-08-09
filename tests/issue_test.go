package tests

import (
	"encoding/binary"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/util"
	"net"
	"testing"
)

func TestInterface(t *testing.T) {

	p := peer.NewPeer("tcp.Acceptor").(cellnet.TCPAcceptor)
	t.Log(p != nil)
}

const maxPacketAddress = "127.0.0.1:16811"

func TestCrackSizePacket(t *testing.T) {
	queue := cellnet.NewEventQueue()

	peerIns := peer.NewGenericPeer("tcp.Acceptor", "server", maxPacketAddress, queue)

	// 设置最大封包约束，默认不约束
	peerIns.(cellnet.TCPSocketOption).SetMaxPacketSize(1000)

	proc.BindProcessorHandler(peerIns, "tcp.ltv", nil)

	peerIns.Start()

	queue.StartLoop()

	conn, err := net.Dial("tcp", maxPacketAddress)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	payload := []byte("hello")

	payloadSize := len(payload)

	fakePayloadSize := payloadSize - 1

	pkt := make([]byte, 2+2+fakePayloadSize)

	// Length, 构建很大的封包
	binary.LittleEndian.PutUint16(pkt, 60000)

	// Type
	binary.LittleEndian.PutUint16(pkt[2:], uint16(1))

	// Value
	copy(pkt[2+2:], payload)

	util.WriteFull(conn, pkt)

	var endBuffer []byte
	_, err = conn.Read(endBuffer)
	if !util.IsEOFOrNetReadError(err) {

		t.Error(err)
		t.FailNow()
	}
}

func TestAutoAllocPort(t *testing.T) {

	p := peer.NewGenericPeer("tcp.Acceptor", "autoacc", ":0", nil)
	p.Start()

	t.Log("auto alloc port:", p.(cellnet.TCPAcceptor).ListenPort())
}
