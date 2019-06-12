package db

import (
	"database/sql"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/peer/mysql"
	"testing"
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
