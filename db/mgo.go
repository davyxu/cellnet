/*
本包只做mongodb的异步操作实现
只实现常用功能

复杂操作可以使用mongodb的原始接口进行操作

*/
package db

import (
	"gopkg.in/mgo.v2"
)

type Config struct {
	URL       string // username:password@tcp(address:port)/dbname   ?authSource=admin
	ConnCount int32  // 连接量
}

type MongoDriver struct {
	*Config
	sesChan chan *mgo.Session
}

func (self *MongoDriver) Start(cfg *Config) error {
	self.Config = cfg

	if self.ConnCount == 0 {
		log.Warnln("DB connection zero")
		return nil
	}

	log.Infof("DB connection %d", self.ConnCount)

	ses, err := mgo.Dial(self.URL)
	if err != nil {

		log.Errorf("%s", err)

		return err
	}
	ses.SetMode(mgo.Strong, true)

	self.sesChan = make(chan *mgo.Session, self.ConnCount)
	self.sesChan <- ses

	// Dial出来的ses底层共享连接池, Copy可以使用这些连接

	for i := 0; i < int(self.ConnCount)-1; i++ {

		self.sesChan <- ses.Copy()
	}

	return nil
}

func (self *MongoDriver) Stop() {

	for i := 0; i < cap(self.sesChan); i++ {
		ses := <-self.sesChan
		ses.Close()
	}

	close(self.sesChan)

}

func (self *MongoDriver) Execute(dbfunc func(sess *mgo.Session)) {

	// 不阻塞当前逻辑
	go func() {

		// 取一个连接
		ses := <-self.sesChan

		defer func() {
			self.sesChan <- ses
		}()

		// 刷新socket
		ses.Refresh()

		// db处理
		dbfunc(ses)

	}()

}

func NewMongoDriver() *MongoDriver {
	return &MongoDriver{}
}
