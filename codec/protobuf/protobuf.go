package protoplus

import (
	cellcodec "github.com/davyxu/cellnet/codec"
	xframe "github.com/davyxu/x/frame"
	"github.com/golang/protobuf/proto"
)

type protobuf struct {
}

func (self *protobuf) Name() string {
	return "protobuf"
}

func (self *protobuf) Encode(msgObj interface{}, ps *xframe.PropertySet) (data interface{}, err error) {

	return proto.Marshal(msgObj.(proto.Message))

}

func (self *protobuf) Decode(data interface{}, msgObj interface{}) error {

	return proto.Unmarshal(data.([]byte), msgObj.(proto.Message))
}

func init() {
	cellcodec.Register(new(protobuf))
}
