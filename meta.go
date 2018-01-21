package cellnet

import (
	"bytes"
	"fmt"
	"path"
	"reflect"
)

// 消息元信息
type MessageMeta struct {
	Name  string       // 消息名称
	Type  reflect.Type // 消息类型
	ID    int          // 消息ID
	Codec Codec        // 消息用到的编码
}

func (self *MessageMeta) NewType() interface{} {
	return reflect.New(self.Type).Interface()
}

var (
	// 消息元信息与消息名称，消息ID和消息类型的关联关系
	metaByName = map[string]*MessageMeta{}
	metaByID   = map[int]*MessageMeta{}
	metaByType = map[reflect.Type]*MessageMeta{}
)

// 注册消息元信息
func RegisterMessageMeta(meta *MessageMeta) {

	if _, ok := metaByName[meta.Name]; ok {
		panic("duplicate message meta register by name: " + meta.Name)
	}

	if _, ok := metaByID[meta.ID]; ok {
		panic(fmt.Sprintf("duplicate message meta register by id: %d", meta.ID))
	}

	if _, ok := metaByType[meta.Type]; ok {
		panic(fmt.Sprintf("duplicate message meta register by type: %d", meta.ID))
	}

	metaByName[meta.Name] = meta
	metaByID[meta.ID] = meta
	metaByType[meta.Type] = meta

}

// 根据名字查找消息元信息
func MessageMetaByName(name string) *MessageMeta {
	if v, ok := metaByName[name]; ok {
		return v
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

// 消息全名
func MessageFullName(rtype reflect.Type) string {

	if rtype == nil {
		panic("empty msg type")
	}

	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	var b bytes.Buffer
	b.WriteString(path.Base(rtype.PkgPath()))
	b.WriteString(".")
	b.WriteString(rtype.Name())

	return b.String()

}

// 根据id查找消息元信息
func MessageMetaByID(id int) *MessageMeta {
	if v, ok := metaByID[id]; ok {
		return v
	}

	return nil
}

func MessageName(msg interface{}) string {

	meta := MessageMetaByType(reflect.TypeOf(msg).Elem())
	if meta == nil {
		return ""
	}

	return meta.Name
}

func MessageID(msg interface{}) int {

	meta := MessageMetaByType(reflect.TypeOf(msg).Elem())
	if meta == nil {
		return 0
	}

	return int(meta.ID)
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
