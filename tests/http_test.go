package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/json"
	_ "github.com/davyxu/cellnet/comm/httppeer"
	"reflect"
	"testing"
)

//
//func TestHttp(t *testing.T) {
//
//	p := cellnet.NewPeer("http.Acceptor")
//	infoSetter := p.(cellnet.PeerInfo)
//	infoSetter.SetAddress("127.0.0.1:8080")
//
//	p.Start()
//
//	queue := cellnet.NewEventQueue()
//	queue.StartLoop()
//	queue.Wait()
//}

type HttpEchoACK struct {
	Msg   string
	Value int32
}

func (self *HttpEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: cellnet.MustGetCodec("json"),
		Name:  "/hello",
		Type:  reflect.TypeOf((*HttpEchoACK)(nil)).Elem(),
	})

}
