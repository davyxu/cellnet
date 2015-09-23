package nexus

import (
	"github.com/davyxu/cellnet"
)

type configDefine struct {
	Listen string
	Join   string
}

var config configDefine

func init() {

	//cellnet.ReadConfig(&config)

}
