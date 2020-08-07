package cellnet

import (
	"path"
	"reflect"
	"sync"
)

// 消息元信息
type MessageMeta struct {
	Codec Codec        // 消息用到的编码
	Type  reflect.Type // 消息类型, 注册时使用指针类型

	New func() interface{} // 直接生成

	ID int // 消息ID (二进制协议中使用)

	ctxListGuard sync.RWMutex
	ctxList      []*context

	name     string
	fullName string // 包+.+消息名
}

type context struct {
	name string
	data interface{}
}

// 注意, 这里的包名, 取的是文件夹名, 而不是代码中的声明名, 尽量让两者统一
func (self *MessageMeta) init() {
	self.name = self.Type.Name()
	self.fullName = path.Base(self.Type.PkgPath()) + "." + self.Type.Name()
}

func (self *MessageMeta) TypeName() string {

	if self == nil {
		return ""
	}

	return self.name
}

func (self *MessageMeta) FullName() string {

	if self == nil {
		return ""
	}

	return self.fullName
}

// 创建meta类型的实例
func (self *MessageMeta) NewType() interface{} {

	if self.New != nil {
		return self.New()
	}

	if self.Type == nil {
		return nil
	}

	return reflect.New(self.Type).Interface()
}

// 为meta对应的名字绑定上下文
func (self *MessageMeta) SetContext(name string, data interface{}) *MessageMeta {

	self.ctxListGuard.Lock()
	defer self.ctxListGuard.Unlock()

	for _, ctx := range self.ctxList {

		if ctx.name == name {
			ctx.data = data
			return self
		}
	}

	self.ctxList = append(self.ctxList, &context{
		name: name,
		data: data,
	})

	return self
}

// 获取meta对应的名字绑定上下文
func (self *MessageMeta) GetContext(key string) (interface{}, bool) {

	self.ctxListGuard.RLock()
	defer self.ctxListGuard.RUnlock()

	for _, ctx := range self.ctxList {

		if ctx.name == key {
			return ctx.data, true
		}
	}

	return nil, false
}

// 按字符串格式取context
func (self *MessageMeta) GetContextAsString(key, defaultValue string) string {

	if v, ok := self.GetContext(key); ok {

		if str, ok := v.(string); ok {
			return str
		}
	}

	return defaultValue
}

// 按字符串格式取context
func (self *MessageMeta) GetContextAsInt(name string, defaultValue int) int {

	if v, ok := self.GetContext(name); ok {

		if intV, ok := v.(int); ok {
			return intV
		}
	}

	return defaultValue
}
