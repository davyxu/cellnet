package cellnet

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var (
	// 本进程的ID
	RegionID int32
)

type configDefine struct {
	RegionID int32
	CellLog  bool
}

var config configDefine

func init() {

	ReadConfig(&config)

	RegionID = config.RegionID

	if config.CellLog {
		log.Printf("[cellnet] Region: %d", RegionID)
	}

}

func ReadConfig(data interface{}) {

	if len(os.Args) < 2 {
		return
	}

	if _, err := toml.DecodeFile(os.Args[1], data); err != nil {
		log.Println(err)
		return
	}

}
