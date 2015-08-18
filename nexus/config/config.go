package config

import (
	"flag"
	"github.com/davyxu/cellnet"
)

type NexusConfig struct {
	ListenAddress string
	JoinAddress   string
	TestCase      string // 单元测试时使用的测试用例名称
}

var Data NexusConfig

func init() {

	flag.StringVar(&Data.ListenAddress, "listen", "", "listen tcp address")

	flag.StringVar(&Data.JoinAddress, "join", "", "join to address")

	flag.StringVar(&Data.TestCase, "test", "", "name of testcase")

	var regionID int
	flag.IntVar(&regionID, "region", 0, "region id")

	flag.Parse()

	// 通知底层, 本进程的regionid是多少
	cellnet.InitRegionID(int32(regionID))

}
