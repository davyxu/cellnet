package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/peer/tcp" // 注册TCP Peer
	_ "github.com/davyxu/cellnet/proc/tcp" // 注册TCP Processor
	"github.com/davyxu/cellnet/util"
	"reflect"
)

const peerAddress = "127.0.0.1:17701"

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

// 将消息注册到系统
func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    int(util.StringHash("main.TestEchoACK")),
	})
}

var clientmode = flag.Int("clientmode", 0, "0: for async recv, 1: for async rpc, 2: for sync rpc")

func main() {

	flag.Parse()

	server()

	switch *clientmode {
	case 0:
		fmt.Println("client mode: async callback")
		clientAsyncCallback()
	case 1:
		fmt.Println("client mode: async rpc")
		clientAsyncRPC()
	case 2:
		fmt.Println("client mode: sync rpc")
		clientSyncRPC()
	}

}
