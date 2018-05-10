package rpc

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/binary"
	"github.com/davyxu/cellnet/util"
	"reflect"
)

type RemoteCallREQ struct {
	MsgID  uint16
	Data   []byte
	CallID int64
}

type RemoteCallACK struct {
	MsgID  uint16
	Data   []byte
	CallID int64
}

func (self *RemoteCallREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *RemoteCallACK) String() string { return fmt.Sprintf("%+v", *self) }

func (self *RemoteCallREQ) GetMsgID() uint16   { return self.MsgID }
func (self *RemoteCallREQ) GetMsgData() []byte { return self.Data }
func (self *RemoteCallREQ) GetCallID() int64   { return self.CallID }
func (self *RemoteCallACK) GetMsgID() uint16   { return self.MsgID }
func (self *RemoteCallACK) GetMsgData() []byte { return self.Data }
func (self *RemoteCallACK) GetCallID() int64   { return self.CallID }

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*RemoteCallREQ)(nil)).Elem(),
		ID:    int(util.StringHash("rpc.RemoteCallREQ")),
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*RemoteCallACK)(nil)).Elem(),
		ID:    int(util.StringHash("rpc.RemoteCallACK")),
	})
}
