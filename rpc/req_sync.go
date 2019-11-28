package rpc

import (
	"time"
)

// 同步RPC请求, ud: peer/session,   reqMsg:请求用的消息, 返回消息为返回值
func CallSync(ud interface{}, reqMsg interface{}, timeout time.Duration) (interface{}, error) {

	ses, err := getPeerSession(ud)

	if err != nil {
		return nil, err
	}

	ret := make(chan interface{})
	// 发送RPC请求
	req := createRequest(func(feedbackMsg interface{}) {
		ret <- feedbackMsg
	})

	req.Send(ses, reqMsg)

	// 等待RPC回复
	select {
	case v := <-ret:
		return v, nil
	case <-time.After(timeout):

		// 清理请求
		getRequest(req.id)

		return nil, ErrTimeout
	}
}
