package cellnet

import (
	"fmt"
	"reflect"
)

type ContentReflector interface {
	Reflect(interface{}) string
}

var contentReflector ContentReflector

// 设置内容解析器
func SetContentReflector(r ContentReflector) {
	contentReflector = r
}

// 获取内容解析器
func GetContentReflector() ContentReflector {
	return contentReflector
}

func ReflectContent(d interface{}) string {

	if contentReflector != nil {
		return contentReflector.Reflect(d)
	}

	return fmt.Sprintf("%s|%v", reflectName(d), d)

}

// 反射获得名称传入的内容
func reflectName(d interface{}) string {

	if d == nil {
		return ""
	}

	t := reflect.TypeOf(d)

	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}

	return t.Name()
}
