package rpc

import (
	"errors"
	"github.com/davyxu/cellnet"
	"time"
)

var ErrTimeout = errors.New("time out")

// 发出请求, 接收到服务器返回后才返回, ud: peer/session,   reqMsg:请求用的消息, ackMsgName: 返回消息类型名, 返回消息为返回值
func CallSync(ud interface{}, reqMsg interface{}, ackMsgName string, timeout time.Duration) (interface{}, error) {

	ses, p, err := getPeerSession(ud)

	if err != nil {
		return nil, err
	}

	ret := make(chan interface{})

	rpcid, err := buildRecvHandler(p, ackMsgName, NewRetChanHandler(ret))
	if err != nil {
		return nil, err
	}

	// 发送RPC请求
	ev := cellnet.NewEvent(cellnet.Event_Send, ses)
	ev.TransmitTag = rpcid
	ev.Msg = reqMsg
	ev.ChainSend = ChainSend()
	ses.RawSend(ev)

	// 等待RPC回复
	select {
	case v := <-ret:
		return v, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}
