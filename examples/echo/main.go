package main

import (
	"flag"
	"fmt"
	_ "github.com/davyxu/cellnet/peer/tcp" // 注册TCP Peer
	_ "github.com/davyxu/cellnet/proc/tcp" // 注册TCP Processor
)

const peerAddress = "127.0.0.1:17701"

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
