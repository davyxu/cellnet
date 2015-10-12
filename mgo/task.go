package mgo

import (
	"github.com/davyxu/cellnet"
	"gopkg.in/mgo.v2"
	"reflect"
)

func (self *mongoDriver) readTask(evq cellnet.EventQueue, usercallback interface{}, execCallback func(*mgo.Database, interface{}) error) {

	callbackType := reflect.TypeOf(usercallback)
	callbackValue := reflect.ValueOf(usercallback)

	resultPtrType := callbackType.In(0)
	// 取回调函数第一个参数类型
	resultType := resultPtrType.Elem()

	// 实例化
	result := reflect.New(resultType)

	db := self.fetch()

	err := execCallback(db, result.Interface())

	self.back(db)

	//将回调函数投递到主线程
	evq.PostData(func() {

		errValue := reflect.ValueOf(&err).Elem()

		// 错误时, 返回空
		if err == nil {

			// http://play.golang.org/p/TZyOLzu2y-
			callbackValue.Call([]reflect.Value{result, errValue})
		} else {
			// http://grokbase.com/t/gg/golang-nuts/13adpw445j/go-nuts-reflect-set-to-nil
			callbackValue.Call([]reflect.Value{reflect.Zero(resultPtrType), errValue})
		}

	})
}

func (self *mongoDriver) writeTask(evq cellnet.EventQueue, usercallback func(error), execCallback func(*mgo.Database) error) {

	db := self.fetch()

	err := execCallback(db)

	self.back(db)

	if usercallback != nil {
		//将回调函数投递到主线程
		evq.PostData(func() {

			usercallback(err)

		})
	}

}
