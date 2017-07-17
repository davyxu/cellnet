package websocket

import "bytes"

// 格式:  消息名\n+json文本

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
