package protoplus

import (
	cellcodec "github.com/davyxu/cellnet/codec"
	"github.com/davyxu/protoplus/api/golang"
	"github.com/davyxu/protoplus/api/golang/wire"
	"github.com/davyxu/x/container"
)

type protoplus struct {
}

func (self *protoplus) Name() string {
	return "protoplus"
}

func (self *protoplus) Encode(msgObj any, ps *xcontainer.Mapper) (data any, err error) {

	return ppgo.Marshal(msgObj.(ppgo.Struct))

}

func (self *protoplus) Decode(data any, msgObj any) error {

	return ppgo.Unmarshal(data.([]byte), msgObj.(wire.Struct))
}

func init() {
	cellcodec.Register(new(protoplus))
}
