package gorillaws

import (
	"bytes"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/gorilla/websocket"
)

func parsePacket(pkt []byte) (msgName string, data []byte) {

	for index, d := range pkt {

		if d == '\n' {
			msgName = string(pkt[:index])
			data = pkt[index+1:]
			return
		}
	}

	return
}

func composePacket(msgName string, data []byte) []byte {

	var b bytes.Buffer
	b.WriteString(msgName)
	b.WriteString("\n")
	b.Write(data)

	return b.Bytes()
}

type WSMessageTransmitter struct {
}

func (WSMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	conn, ok := ses.Raw().(*websocket.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil, nil
	}

	t, raw, err := conn.ReadMessage()

	switch t {
	case websocket.TextMessage:
		msgName, userPacket := parsePacket(raw)

		if msgName != "" {

			meta := cellnet.MessageMetaByFullName(msgName)

			if meta == nil || meta.Codec == nil {
				return nil, cellnet.NewErrorContext("codec error", msgName)
			}

			msg, _, err = codec.DecodeMessage(meta.ID, userPacket)
		}

	case websocket.BinaryMessage:
		// TODO 实现二进制
	}

	return
}

func (WSMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	conn, ok := ses.Raw().(*websocket.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil
	}

	data, meta, err := codec.EncodeMessage(msg, nil)
	if err != nil {
		return err
	}

	// 组websocket包
	raw := composePacket(meta.FullName(), data)

	conn.WriteMessage(websocket.TextMessage, raw)

	return nil
}
