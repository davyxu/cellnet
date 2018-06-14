package http

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"io"
	"io/ioutil"
	"net/http"
)

type MessageRespond struct {
	StatusCode int
	Msg        interface{}
	CodecName  string
}

func (self *MessageRespond) String() string {
	return fmt.Sprintf("Code: %d Msg: %+v CodeName: %s", self.StatusCode, self.Msg, self.CodecName)
}

func (self *MessageRespond) WriteRespond(ses *httpSession) error {
	peerInfo := ses.Peer().(cellnet.PeerProperty)

	httpCodec := codec.GetCodec(self.CodecName)

	if httpCodec == nil {
		return errors.New("ResponseCodec not found:" + self.CodecName)
	}

	msg := self.Msg

	log.Debugf("#http.send(%s) '%s' %s | [%d] Message(%s) %s",
		peerInfo.Name(),
		ses.req.Method,
		ses.req.URL.Path,
		self.StatusCode,
		cellnet.MessageToName(msg),
		cellnet.MessageToString(msg))

	// 将消息编码为字节数组
	var data interface{}
	data, err := httpCodec.Encode(msg, nil)

	if err != nil {
		return err
	}

	ses.resp.Header().Set("Content-Type", httpCodec.MimeType()+";charset=UTF-8")
	ses.resp.WriteHeader(http.StatusOK)

	bodyData, err := ioutil.ReadAll(data.(io.Reader))
	if err != nil {
		return err
	}

	ses.resp.Write(bodyData)

	return nil
}
