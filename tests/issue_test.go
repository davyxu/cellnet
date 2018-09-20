package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"sync"
	"testing"
)

func TestContextSet(t *testing.T) {

	p := peer.NewPeer("tcp.Acceptor").(cellnet.TCPAcceptor)
	p.(cellnet.ContextSet).SetContext("sd", nil)

	if v, ok := p.(cellnet.ContextSet).GetContext("sd"); ok && v == nil {

	} else {
		t.FailNow()
	}

	var connMap = new(sync.Map)
	if p.(cellnet.ContextSet).FetchContext("sd", &connMap) && connMap == nil {

	} else {
		t.FailNow()
	}
}

func TestAutoAllocPort(t *testing.T) {

	p := peer.NewGenericPeer("tcp.Acceptor", "autoacc", ":0", nil)
	p.Start()

	t.Log("auto alloc port:", p.(cellnet.TCPAcceptor).Port())
}
