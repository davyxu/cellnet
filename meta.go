package cellnet

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strings"
)

// 消息元信息
type MessageMeta struct {
	Codec Codec        // 消息用到的编码
	Type  reflect.Type // 消息类型

	ID int // 消息ID (二进制协议中使用)
}

func (self *MessageMeta) TypeName() string {

	if self == nil {
		return ""
	}

	if self.Type.Kind() == reflect.Ptr {
		return self.Type.Elem().Name()
	}

	return self.Type.Name()
}

func (self *MessageMeta) FullName() string {

	if self == nil {
		return ""
	}

	rtype := self.Type
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	var sb strings.Builder
	sb.WriteString(path.Base(rtype.PkgPath()))
	sb.WriteString(".")
	sb.WriteString(rtype.Name())

	return sb.String()
}

func (self *MessageMeta) NewType() interface{} {
	return reflect.New(self.Type).Interface()
}

var (
	// 消息元信息与消息名称，消息ID和消息类型的关联关系
	metaByFullName = map[string]*MessageMeta{}
	metaByID       = map[int]*MessageMeta{}
	metaByType     = map[reflect.Type]*MessageMeta{}
)

/*
http消息
Method URL -> Meta
Type -> Meta

非http消息
ID -> Meta
Type -> Meta

*/

// 注册消息元信息
func RegisterMessageMeta(meta *MessageMeta) {

	// 非http类,才需要包装Type必须唯一

	if _, ok := metaByType[meta.Type]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by type: %d", meta.ID))
	} else {
		metaByType[meta.Type] = meta
	}

	if _, ok := metaByFullName[meta.FullName()]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by fullname: %s", meta.FullName()))
	} else {
		metaByFullName[meta.FullName()] = meta
	}

	if meta.ID == 0 {
		panic("message meta require 'ID' field: " + meta.TypeName())
	}

	if _, ok := metaByID[meta.ID]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by id: %d", meta.ID))
	} else {
		metaByID[meta.ID] = meta
	}

}

// 根据名字查找消息元信息
func MessageMetaByFullName(name string) *MessageMeta {
	if v, ok := metaByFullName[name]; ok {
		return v
	}

	return nil
}

func MessageMetaVisit(nameRule string, callback func(meta *MessageMeta) bool) error {
	exp, err := regexp.Compile(nameRule)
	if err != nil {
		return err
	}

	for name, meta := range metaByFullName {
		if exp.MatchString(name) {

			if !callback(meta) {
				return nil
			}

		}
	}

	return nil
}

// 根据类型查找消息元信息
func MessageMetaByType(t reflect.Type) *MessageMeta {

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if v, ok := metaByType[t]; ok {
		return v
	}

	return nil
}

// 根据id查找消息元信息
func MessageMetaByID(id int) *MessageMeta {
	if v, ok := metaByID[id]; ok {
		return v
	}

	return nil
}

func MessageToName(msg interface{}) string {

	if msg == nil {
		return ""
	}

	meta := MessageMetaByType(indirectType(msg))
	if meta == nil {
		return ""
	}

	return meta.TypeName()
}

func indirectType(msg interface{}) reflect.Type {
	t := reflect.TypeOf(msg)
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}

	return t
}

func MessageToID(msg interface{}) int {

	if msg == nil {
		return 0
	}

	reflect.TypeOf(msg)

	meta := MessageMetaByType(indirectType(msg))
	if meta == nil {
		return 0
	}

	return int(meta.ID)
}

func MessageSize(msg interface{}) int {

	if msg == nil {
		return 0
	}

	// 获取消息元信息
	meta := MessageMetaByType(reflect.TypeOf(msg))
	if meta == nil {
		return 0
	}

	// 将消息编码为字节数组
	raw, err := meta.Codec.Encode(msg)

	if err != nil {
		return 0
	}

	return len(raw.([]byte))
}

func MessageToString(msg interface{}) string {

	if msg == nil {
		return ""
	}

	if stringer, ok := msg.(interface {
		String() string
	}); ok {
		return stringer.String()
	}

	return ""
}
