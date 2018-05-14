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
