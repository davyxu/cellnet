package natstransport

import (
	cellcodec "github.com/davyxu/cellnet/codec"
	xbytes "github.com/davyxu/x/bytes"
)

func SendMessage(msg any) ([]byte, error) {
	data, meta, err := cellcodec.Encode(msg, nil)
	if err != nil {
		return nil, err
	}
	payload := make([]byte, len(data)+2)
	writer := xbytes.NewWriter(payload)
	writer.WriteUint32(uint32(meta.Id))
	writer.Write(data)
	return payload, err
}

func RecvMessage(payload []byte) (any, error) {
	reader := xbytes.NewReader(payload)
	msgID, err := reader.ReadUint32()
	if err != nil {
		return nil, err
	}

	msgData := payload[reader.Size():]

	msg, _, err := cellcodec.Decode(int(msgID), msgData)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
