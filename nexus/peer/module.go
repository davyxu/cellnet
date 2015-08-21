package peer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/dispatcher"
	"github.com/davyxu/cellnet/nexus/config"
)

var Dispatcher = dispatcher.NewPacketDispatcher()

func init() {
	cellnet.RegisterModuleEntry(func() {

		listenNexus()

		joinAddr := config.Data.JoinAddress

		if joinAddr != "" {

			joinNexus(joinAddr)
		}

	})
}
