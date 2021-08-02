package protoplus

import (
	cellcodec "github.com/davyxu/cellnet/codec"
	"github.com/davyxu/protoplus/api/golang"
	"github.com/davyxu/protoplus/api/golang/wire"
	xframe "github.com/davyxu/x/frame"
)

type protoplus struct {
}

func (self *protoplus) Name() string {
	return "protoplus"
}

func (self *protoplus) Encode(msgObj interface{}, ps *xframe.Mapper) (data interface{}, err error) {

	return ppgo.Marshal(msgObj.(ppgo.Struct))

}

func (self *protoplus) Decode(data interface{}, msgObj interface{}) error {

	return ppgo.Unmarshal(data.([]byte), msgObj.(wire.Struct))
}

func init() {
	cellcodec.Register(new(protoplus))
}
