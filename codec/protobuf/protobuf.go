package protoplus

import (
	cellcodec "github.com/davyxu/cellnet/codec"
	xframe "github.com/davyxu/x/frame"
	pb "github.com/golang/protobuf/proto"
)

type protobuf struct {
}

func (self *protobuf) Name() string {
	return "protobuf"
}

func (self *protobuf) Encode(msgObj interface{}, ps *xframe.PropertySet) (data interface{}, err error) {

	return pb.Marshal(msgObj.(pb.Message))

}

func (self *protobuf) Decode(data interface{}, msgObj interface{}) error {

	return pb.Unmarshal(data.([]byte), msgObj.(pb.Message))
}

func init() {
	cellcodec.Register(new(protobuf))
}
