package redix

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"sync/atomic"
	"time"
)

type redisConnector struct {
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRedisParameter

	*pool.Pool

	readyValue int64
}

func (self *redisConnector) IsReady() bool {
	return atomic.LoadInt64(&self.readyValue) != 0
}

func (self *redisConnector) TypeName() string {
	return "redix.Connector"
}

func (self *redisConnector) Raw() interface{} {
	return self.Pool
}

func (self *redisConnector) Operate(callback func(client interface{}) interface{}) interface{} {

	c, err := self.Pool.Get()
	if err != nil {
		log.Errorf("get client failed, %s", err)
		return err
	}

	defer func() {
		self.Pool.Put(c)
	}()

	return callback(c)
}

func (self *redisConnector) tryConnect() {
	var err error

	for {
		self.Pool, err = pool.NewCustom("tcp", self.Address(), self.PoolConnCount, func(network, addr string) (*redis.Client, error) {
			client, err := redis.Dial(network, addr)
			if err != nil {
				return nil, err
			}
			if len(self.Password) > 0 {
				if err = client.Cmd("AUTH", self.Password).Err; err != nil {
					client.Close()
					return nil, err
				}
			}
			if self.DBIndex > 0 {
				if err = client.Cmd("SELECT", self.DBIndex).Err; err != nil {
					client.Close()
					return nil, err
				}
			}

			log.Infof("Create redis pool connection: %s | %s", addr, util.StackToString(10))

			return client, nil
		})

		if err != nil {
			log.Errorln("Redis connect failed:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		break
	}

	atomic.StoreInt64(&self.readyValue, 1)
}

func (self *redisConnector) Start() cellnet.Peer {

	go self.tryConnect()

	return self
}

func (self *redisConnector) Stop() {

}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {

		self := &redisConnector{}
		self.CoreRedisParameter.Init()

		return self
	})
}
