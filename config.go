package cellnet

import (
	"log"
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

	RegionID = config.RegionID

	if config.CellLog {
		log.Printf("[cellnet] Region: %d", RegionID)
	}

}

func EnableLog(v bool) {
	config.CellLog = v
}
