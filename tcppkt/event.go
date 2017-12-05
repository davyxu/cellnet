package tcppkt

import "github.com/davyxu/cellnet"

type RecvMsgEvent struct {
	Ses cellnet.Session
	Msg interface{}
}

type SendMsgEvent struct {
	Ses cellnet.Session
	Msg interface{}
}
