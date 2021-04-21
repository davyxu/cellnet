package cellmeta

import (
	"github.com/davyxu/cellnet"
	xframe "github.com/davyxu/x/frame"
	"path"
	"reflect"
)

// 消息元信息
type Meta struct {
	xframe.PropertySet

	Type  reflect.Type       // 消息类型, 注册时使用指针类型
	Codec cellnet.Codec      // 消息用到的编码
	New   func() interface{} // 直接生成

	ID       int // 消息ID (二进制协议中使用)
	name     string
	fullName string // 包+.+消息名

}

// 注意, 这里的包名, 取的是文件夹名, 而不是代码中的声明名, 尽量让两者统一
func (self *Meta) init() {
	self.name = self.Type.Name()
	self.fullName = path.Base(self.Type.PkgPath()) + "." + self.Type.Name()
}

func (self *Meta) TypeName() string {

	if self == nil {
		return ""
	}

	return self.name
}

func (self *Meta) FullName() string {

	if self == nil {
		return ""
	}

	return self.fullName
}

// 创建meta类型的实例
func (self *Meta) NewType() interface{} {

	if self.New != nil {
		return self.New()
	}

	if self.Type == nil {
		return nil
	}

	return reflect.New(self.Type).Interface()
}
