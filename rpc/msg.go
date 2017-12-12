package rpc

import (
	"fmt"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/binary"
	"reflect"
)

type RemoteCallMsg interface {
	GetMsgID() uint32
	GetMsgData() []byte
	GetCallID() int64
}

type RemoteCallREQ struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}

func (self *RemoteCallREQ) GetMsgID() uint32 {
	return self.MsgID
}

func (self *RemoteCallREQ) GetMsgData() []byte {
	return self.Data
}

func (self *RemoteCallREQ) GetCallID() int64 {
	return self.CallID
}

type RemoteCallACK struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}

func (self *RemoteCallACK) GetMsgID() uint32 {
	return self.MsgID
}

func (self *RemoteCallACK) GetMsgData() []byte {
	return self.Data
}

func (self *RemoteCallACK) GetCallID() int64 {
	return self.CallID
}

func (self *RemoteCallREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *RemoteCallACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta("binary", "rpc.RemoteCallREQ", reflect.TypeOf((*RemoteCallREQ)(nil)).Elem(), 63101)
	cellnet.RegisterMessageMeta("binary", "rpc.RemoteCallACK", reflect.TypeOf((*RemoteCallACK)(nil)).Elem(), 63102)
}
