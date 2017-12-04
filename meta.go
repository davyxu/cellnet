package cellnet

import (
	"bytes"
	"fmt"
	"path"
	"reflect"
)

// 消息元信息
type MessageMeta struct {
	Type  reflect.Type // 消息类型
	Name  string       // 消息名称
	ID    uint32       // 消息ID
	Codec Codec        // 消息用到的编码
}

var (
	// 消息元信息与消息名称，消息ID和消息类型的关联关系
	metaByName = map[string]*MessageMeta{}
	metaByID   = map[uint32]*MessageMeta{}
	metaByType = map[reflect.Type]*MessageMeta{}
)

// 注册消息元信息
func RegisterMessageMeta(codecName string, name string, msgType reflect.Type, id uint32) {

	meta := &MessageMeta{
		Type:  msgType,
		Name:  name,
		ID:    id,
		Codec: FetchCodec(codecName),
	}

	if meta.Codec == nil {
		panic("codec not register! " + codecName)
	}

	if _, ok := metaByName[name]; ok {
		panic("duplicate message meta register by name: " + name)
	}

	if _, ok := metaByID[meta.ID]; ok {
		panic(fmt.Sprintf("duplicate message meta register by id: %d", meta.ID))
	}

	if _, ok := metaByType[msgType]; ok {
		panic(fmt.Sprintf("duplicate message meta register by type: %d", meta.ID))
	}

	metaByName[name] = meta
	metaByID[meta.ID] = meta
	metaByType[msgType] = meta
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
func MessageMetaByID(id uint32) *MessageMeta {
	if v, ok := metaByID[id]; ok {
		return v
	}

	return nil
}
