package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/mgo"
	"github.com/davyxu/cellnet/test"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
)

var signal *test.SignalTester

type char struct {
	Name string
	HP   int
}

func db() {
	pipe := cellnet.NewEventPipe()

	evq := pipe.AddQueue()

	pipe.Start()

	db := mgo.NewDB()

	var err error

	err = db.Start(&mgo.Config{
		URL:     "127.0.0.1:27017/test",
		ShowLog: true,
	})

	if err != nil {
		log.Println("db connect failed:", err)
		signal.Fail()
		return
	}

	db.FindOne(evq, "test", bson.M{"name": "davy"}, func(c *char, _ error) {

		// 没有记录, 创建
		if c == nil {
			db.Insert(evq, "test", &char{Name: "davy", HP: 10}, nil)
			db.Insert(evq, "test", &char{Name: "zerg", HP: 90}, func(_ error) {

				db.FindOne(evq, "test", bson.M{"name": "davy"}, func(c *char, _ error) {

					if c == nil {
						signal.Log("can not found record")
						signal.Fail()
					} else {
						log.Println(c)
						signal.Done(1)
					}

				})

			})

			// 有记录, 搞定
		} else {

			log.Println(c)
			signal.Done(1)

		}

	})

	db.Update(evq, "test", bson.M{"name": "davy"}, &char{Name: "davy", HP: 1}, func(err error) {

		if err != nil {
			signal.Log("update failed")
			signal.Fail()
		}

		db.FindOne(evq, "test", bson.M{"name": "davy"}, func(c *char, _ error) {

			if c == nil {
				signal.Log("update failed")
				signal.Fail()
			} else {
				if c.HP != 1 {
					signal.Fail()
				} else {
					signal.Done(2)
				}

			}

		})
	})

	signal.WaitAndExpect(1, "find failed")
	signal.WaitAndExpect(2, "update failed")

}

func TestMongoDB(t *testing.T) {

	signal = test.NewSignalTester(t)

	db()
}
