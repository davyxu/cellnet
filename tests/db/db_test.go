package db

import (
	"database/sql"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/peer/mysql"
	_ "github.com/davyxu/cellnet/peer/redix"
	"github.com/mediocregopher/radix.v2/redis"
	"testing"
	"time"
)

func TestMySQL(t *testing.T) {
	// 从地址中选择mysql数据库，这里选择mysql系统库
	p := peer.NewGenericPeer("mysql.Connector", "mysqldemo", "root:123456@(localhost:3306)/mysql", nil)
	p.(cellnet.MySQLConnector).SetConnectionCount(3)

	// 阻塞
	p.Start()

	op := p.(cellnet.MySQLOperator)

	op.Operate(func(rawClient interface{}) interface{} {

		client := rawClient.(*sql.DB)

		// 查找默认用户
		mysql.NewWrapper(client).Query("select User from user").Each(func(wrapper *mysql.Wrapper) bool {

			var name string
			wrapper.Scan(&name)
			fmt.Println(name)
			return true
		})

		return nil
	})

}

func TestRedix(t *testing.T) {

	peer := peer.NewGenericPeer("redix.Connector", "redis", "127.0.0.1:6379", nil)
	peer.Start()

	for i := 0; i < 5; i++ {

		if peer.(cellnet.PeerReadyChecker).IsReady() {
			break
		}

		time.Sleep(time.Millisecond * 200)
	}

	if !peer.(cellnet.PeerReadyChecker).IsReady() {
		t.FailNow()
	}

	op := peer.(cellnet.RedisPoolOperator)

	op.Operate(func(rawClient interface{}) interface{} {

		client := rawClient.(*redis.Client)
		client.Cmd("SET", "mykey", "myvalue")

		v, err := client.Cmd("GET", "mykey").Str()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if v != "myvalue" {
			t.FailNow()
		}

		return nil
	})

}
