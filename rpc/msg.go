package rpc

import (
	"fmt"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/binary"
	"reflect"
)

type RemoteCallREQ struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}

type RemoteCallACK struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}

func (self *RemoteCallREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *RemoteCallACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta("binary", "rpc.RemoteCallREQ", reflect.TypeOf((*RemoteCallREQ)(nil)).Elem(), 11)
	cellnet.RegisterMessageMeta("binary", "rpc.RemoteCallACK", reflect.TypeOf((*RemoteCallACK)(nil)).Elem(), 12)
}
