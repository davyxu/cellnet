package rpc

import (
	"time"
)

// 发出请求, 接收到服务器返回后才返回, ud: peer/session,   reqMsg:请求用的消息, 返回消息为返回值
func CallSync(ud interface{}, reqMsg interface{}, timeout time.Duration) (interface{}, error) {

	ses, _, err := getPeerSession(ud)

	if err != nil {
		return nil, err
	}

	// 发送RPC请求
	req := createRequest(timeout)

	ret := make(chan interface{})

	req.onRecv = func(feedbackMsg interface{}) {
		ret <- feedbackMsg
	}

	sendRequest(ses, reqMsg, req.id)

	// 等待RPC回复
	select {
	case v := <-ret:
		return v, nil
	case <-time.After(req.timeout):
		return nil, ErrTimeout
	}
}
