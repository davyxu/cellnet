package relay

type BroadcasterFunc func(event *RecvMsgEvent) bool

var bcFunc BroadcasterFunc

// 设置广播函数
func SetBroadcaster(callback BroadcasterFunc) {

	bcFunc = callback
}
