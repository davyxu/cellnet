package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/redix"
	"github.com/mediocregopher/radix.v2/redis"
	"testing"
	"time"
)

func TestRedix(t *testing.T) {

	peer := peer.NewGenericPeer("redix.Connector", "redis", "127.0.0.1:16379", nil)
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
