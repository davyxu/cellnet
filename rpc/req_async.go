package rpc

import (
	"time"
)

// ud: peer/session,   reqMsg:请求用的消息, userCallback: 返回消息类型回调 func( ackMsg *ackMsgType)
func Call(sesOrPeer interface{}, reqMsg interface{}, timeout time.Duration, userCallback func(raw interface{})) {

	ses, _, err := getPeerSession(sesOrPeer)

	if err != nil {
		userCallback(err)
		return
	}

	// 发送RPC请求
	req := createRequest(timeout)

	req.onRecv = userCallback

	sendRequest(ses, reqMsg, req.id)

	// 等待RPC回复
	time.AfterFunc(timeout, func() {

		if requestExists(req.id) {
			userCallback(ErrTimeout)
		}
	})
}
