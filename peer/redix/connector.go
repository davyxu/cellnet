package redix

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"sync"
	"time"
)

type redisConnector struct {
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRedisParameter

	pool      *pool.Pool
	poolGuard sync.RWMutex
}

func (self *redisConnector) IsReady() bool {
	return self.Pool() != nil
}

func (self *redisConnector) TypeName() string {
	return "redix.Connector"
}

func (self *redisConnector) Pool() *pool.Pool {
	self.poolGuard.RLock()
	defer self.poolGuard.RUnlock()

	return self.pool
}

func (self *redisConnector) Raw() interface{} {

	return self.Pool()
}

func (self *redisConnector) Operate(callback func(client interface{}) interface{}) interface{} {

	pool := self.Pool()
	c, err := pool.Get()
	if err != nil {
		log.Errorf("get client failed, %s", err)
		return err
	}

	defer func() {
		pool.Put(c)
	}()

	return callback(c)
}

func (self *redisConnector) tryConnect() {

	for {
		pool, err := pool.NewCustom("tcp", self.Address(), self.PoolConnCount, func(network, addr string) (*redis.Client, error) {

			client, err := redis.DialTimeout(network, addr, time.Second*5)
			if err != nil {
				log.Errorf("redis.Dial %s", err.Error())
				return nil, err
			}

			if len(self.Password) > 0 {
				if err = client.Cmd("AUTH", self.Password).Err; err != nil {
					log.Errorf("redis.Auth %s %s", self.Password, err.Error())
					client.Close()
					return nil, err
				}
			}
			if self.DBIndex > 0 {
				if err = client.Cmd("SELECT", self.DBIndex).Err; err != nil {
					log.Errorf("redis.SELECT %d %s", self.DBIndex, err.Error())
					client.Close()
					return nil, err
				}
			}

			log.Infof("Create redis pool connection: %s name: %s index: %d", addr, self.Name(), self.DBIndex)

			return client, nil
		})

		if err != nil {
			log.Errorln("Redis connect failed:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		self.poolGuard.Lock()
		self.pool = pool
		self.poolGuard.Unlock()

		break
	}

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
