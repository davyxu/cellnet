package protoplus

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/protoplus/api/golang"
	"github.com/davyxu/protoplus/api/golang/wire"
)

type protoplus struct {
}

func (self *protoplus) Name() string {
	return "protoplus"
}

func (self *protoplus) MimeType() string {
	return "application/binary"
}

func (self *protoplus) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {

	return ppgo.Marshal(msgObj.(ppgo.Struct))

}

func (self *protoplus) Decode(data interface{}, msgObj interface{}) error {

	return ppgo.Unmarshal(data.([]byte), msgObj.(wire.Struct))
}

func init() {

	codec.RegisterCodec(new(protoplus))
}
