package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"testing"
)

func TestAutoAllocPort(t *testing.T) {

	p := peer.NewGenericPeer("tcp.Acceptor", "autoacc", ":0", nil)
	p.Start()

	t.Log("auto alloc port:", p.(cellnet.TCPAcceptor).Port())
}
