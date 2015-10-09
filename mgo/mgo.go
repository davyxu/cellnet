/*
本包只做mongodb的异步操作实现
只实现常用功能

复杂操作可以使用mongodb的原始接口进行操作

*/
package mgo

import (
	"github.com/davyxu/cellnet"
	"gopkg.in/mgo.v2"
	"log"
)

type Config struct {
	URL       string // username:password@tcp(address:port)/dbname   ?authSource=admin
	ConnCount int    // 连接量
	ShowLog   bool   // 显示操作日志
}

type mongoDriver struct {
	*Config
	dbChan chan *mgo.Database
}

func (self *mongoDriver) Start(rawcfg interface{}) error {
	self.Config = rawcfg.(*Config)

	if self.ConnCount == 0 {
		self.ConnCount = 1
	}

	if self.ShowLog {
		log.Printf("[mgo] db connection %d", self.ConnCount)
	}

	self.dbChan = make(chan *mgo.Database, self.ConnCount)

	for i := 0; i < self.ConnCount; i++ {

		ses, err := mgo.Dial(self.URL)
		if err != nil {

			if self.ShowLog {
				log.Printf("[mgo] %s", err)
			}

			return err
		}

		ses.SetMode(mgo.Monotonic, true)

		// 默认db
		self.dbChan <- ses.DB("")
	}

	return nil
}

func (self *mongoDriver) back(db *mgo.Database) {
	self.dbChan <- db
}

func (self *mongoDriver) fetch() *mgo.Database {

	return <-self.dbChan
}

func (self *mongoDriver) Stop() {

	for i := 0; i < cap(self.dbChan); i++ {
		db := <-self.dbChan
		db.Session.Close()
	}

	close(self.dbChan)

}

func (self *mongoDriver) Insert(evq cellnet.EventQueue, collName string, doc interface{}, callback func(error)) {

	self.writeTask(evq, callback, func(db *mgo.Database) error {

		if self.ShowLog {
			log.Printf("[mgo] insert '%s' %v", collName, doc)
		}

		err := db.C(collName).Insert(doc)

		if err != nil && self.ShowLog {

			log.Printf("[mgo] insert failed, %s", err)

		}

		return err
	})

}

func (self *mongoDriver) FindOne(evq cellnet.EventQueue, collName string, query interface{}, callback interface{}) {

	self.readTask(evq, callback, func(db *mgo.Database, result interface{}) error {

		if self.ShowLog {
			log.Printf("[mgo] findone '%s' query: %s", collName, query)
		}

		err := db.C(collName).Find(query).One(result)

		if err != nil && self.ShowLog {

			log.Printf("[mgo] findone failed, %s", err)
		}

		return err
	})

}

func (self *mongoDriver) Update(evq cellnet.EventQueue, collName string, selector interface{}, doc interface{}, callback func(error)) {

	self.writeTask(evq, callback, func(db *mgo.Database) error {

		if self.ShowLog {
			log.Printf("[mgo] update '%s' sel: %s update: %v", collName, selector, doc)
		}

		err := db.C(collName).Update(selector, doc)

		if err != nil && self.ShowLog {

			log.Printf("[mgo] update failed, %s", err)

		}

		return err
	})
}

func NewDB() cellnet.KVDatabase {
	return &mongoDriver{}
}
