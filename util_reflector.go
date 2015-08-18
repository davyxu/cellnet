package cellnet

import (
	"fmt"
	"reflect"
)

// 反射获得名称传入的内容
func ReflectName(d interface{}) string {

	if d == nil {
		return ""
	}

	t := reflect.TypeOf(d)

	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}

	return t.Name()
}

// 显示传入的类型和内容
func ReflectContent(d interface{}) string {
	return fmt.Sprintf("%s|%v", ReflectName(d), d)
}
