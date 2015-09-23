package socket

import (
//	"github.com/davyxu/cellnet"
)

type configDefine struct {
	SocketLog bool
}

var config configDefine

func EnableLog(v bool) {
	config.SocketLog = v
}

func init() {

	//cellnet.ReadConfig(&config)

}
