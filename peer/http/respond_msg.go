package http

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

func writeMessageRespond(ses *httpSession, msg interface{}) {
	peerInfo := ses.Peer().(cellnet.PeerProperty)

	log.Debugf("#send Response(%s) %s %s | %s",
		peerInfo.Name(),
		ses.req.URL.Path,
		cellnet.MessageToName(msg),
		cellnet.MessageToString(msg))

	// 获取消息元信息
	meta := cellnet.HttpMetaByResponseType(ses.req.Method, reflect.TypeOf(msg))
	if meta == nil {
		log.Errorln("message not found:", msg)
		return
	}

	// 将消息编码为字节数组
	var data interface{}
	data, err := meta.ResponseCodec.Encode(msg)

	if err != nil {
		log.Errorln("message encode error:", err)
		return
	}

	ses.resp.Write(data.([]byte))
}
