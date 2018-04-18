package relay

import "github.com/davyxu/cellnet"

type BroadcasterFunc func(frontendPeer cellnet.Peer, event *RecvMsgEvent)

type broadcastContext struct {
	backendPeer  cellnet.Peer
	frontendPeer cellnet.Peer
	callback     BroadcasterFunc
}

func (self *broadcastContext) invoke(backendPeer cellnet.Peer, ev *RecvMsgEvent) {

	if self.backendPeer == ev.Ses.Peer() {
		self.callback(self.frontendPeer, ev)
	}
}

var (
	bclist []broadcastContext
)

// 设置广播函数, 回调时，按对应Peer/Session所在的队列中调用
func BindBroadcaster(frontendPeer, backendPeer cellnet.Peer, callback BroadcasterFunc) {

	bclist = append(bclist, broadcastContext{
		backendPeer:  backendPeer,
		frontendPeer: frontendPeer,
		callback:     callback,
	})
}

func broadcast(ev *RecvMsgEvent) {

	backendPeer := ev.Ses.Peer()

	for _, ctx := range bclist {

		ctx.invoke(backendPeer, ev)
	}

}
