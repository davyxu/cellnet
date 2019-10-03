package rpc

import (
	"github.com/davyxu/cellnet"
	"time"
)

// 异步RPC请求
// ud: peer/session,   reqMsg:请求用的消息, userCallback: 返回消息类型回调 func( ackMsg *ackMsgType)
func Call(sesOrPeer interface{}, reqMsg interface{}, timeout time.Duration, userCallback func(raw interface{})) {

	ses, err := getPeerSession(sesOrPeer)

	if err != nil {

		cellnet.SessionQueuedCall(ses, func() {
			userCallback(err)
		})

		return
	}

	// 发送RPC请求
	req := createRequest(func(raw interface{}) {
		cellnet.SessionQueuedCall(ses, func() {
			userCallback(raw)
		})
	})

	req.Send(ses, reqMsg)

	// 等待RPC回复
	time.AfterFunc(timeout, func() {

		// 取出请求，如果存在，调用超时
		if getRequest(req.id) != nil {
			cellnet.SessionQueuedCall(ses, func() {
				userCallback(ErrTimeout)
			})
		}
	})
}
