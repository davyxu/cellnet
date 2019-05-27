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

	p := peer.NewGenericPeer("redix.Connector", "redis", "127.0.0.1:6379", nil)
	p.(cellnet.RedisConnector).SetConnectionCount(1)
	p.Start()

	for i := 0; i < 5; i++ {

		if p.(cellnet.PeerReadyChecker).IsReady() {
			break
		}

		time.Sleep(time.Millisecond * 200)
	}

	if !p.(cellnet.PeerReadyChecker).IsReady() {
		t.FailNow()
	}

	op := p.(cellnet.RedisPoolOperator)

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
