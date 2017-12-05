package packet

import "github.com/davyxu/cellnet"

type RecvMsgEvent struct {
	Ses     cellnet.Session
	Msg     interface{}
	MsgID   int
	MsgData []byte
}

type SendMsgEvent struct {
	Ses cellnet.Session
	Msg interface{}
}
