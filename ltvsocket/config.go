package ltvsocket

import (
	"github.com/davyxu/cellnet"
)

type configDefine struct {
	SocketLog bool
}

var config configDefine

func init() {

	cellnet.ReadConfig(&config)

}
