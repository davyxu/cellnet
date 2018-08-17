package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"testing"
)

func TestInterface(t *testing.T) {

	p := peer.NewPeer("tcp.Acceptor").(cellnet.TCPAcceptor)
	t.Log(p != nil)
}

func TestAutoAllocPort(t *testing.T) {

	p := peer.NewGenericPeer("tcp.Acceptor", "autoacc", ":0", nil)
	p.Start()

	t.Log("auto alloc port:", p.(cellnet.TCPAcceptor).Port())
}
