package cellmeta

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/x/container"
	"path"
	"reflect"
)

// 消息元信息
type Meta struct {
	xcontainer.Mapper

	Type  reflect.Type  // 消息类型, 注册时使用指针类型
	Codec cellnet.Codec // 消息用到的编码
	New   func() any    // 直接生成

	Id       int // 消息ID (二进制协议中使用)
	name     string
	FullName string // 包+.+消息名

}

// 注意, 这里的包名, 取的是文件夹名, 而不是代码中的声明名, 尽量让两者统一
func (self *Meta) init() {
	self.name = self.Type.Name()
	if self.FullName == "" {
		self.FullName = path.Base(self.Type.PkgPath()) + "." + self.Type.Name()
	}
}

func (self *Meta) TypeName() string {

	if self == nil {
		return ""
	}

	return self.name
}

// 创建meta类型的实例
func (self *Meta) NewType() any {

	if self.New != nil {
		return self.New()
	}

	if self.Type == nil {
		return nil
	}

	return reflect.New(self.Type).Interface()
}
