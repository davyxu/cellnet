package packet

import (
	"bytes"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/binary"
	xbytes "github.com/davyxu/x/bytes"
	"reflect"
	"testing"
)

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    int(xbytes.StringHash("tests.TestEchoACK")),
	})
}

func TestLNVPacket(t *testing.T) {

	var b bytes.Buffer

	err := SendLenNameValue(&b, nil, &TestEchoACK{
		Msg: "hello",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	msg, err := RecvLenNameValue(&b, 100000)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(msg)

}
