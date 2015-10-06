package socket

import (
	"github.com/davyxu/cellnet"
)

// Peer间的共享数据
type peerProfile struct {
	*cellnet.EvQueue // 实现事件注册和注入
	name             string

	relay bool
}

func (self *peerProfile) SetName(name string) {
	self.name = name
}

func (self *peerProfile) Name() string {
	return self.name
}

func (self *peerProfile) SetRelayMode(relay bool) {
	self.relay = relay
}

func (self *peerProfile) Event() *cellnet.EvQueue {
	return self.EvQueue
}

func newPeerProfile(queue *cellnet.EvQueue) *peerProfile {

	return &peerProfile{EvQueue: queue}
}
