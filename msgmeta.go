package cellnet

import (
	"path"
	"reflect"
)

type MessageMeta struct {
	Type reflect.Type
	Name string
	ID   uint32
}

var (
	metaByName = map[string]*MessageMeta{}
	metaByID   = map[uint32]*MessageMeta{}
)

// 注册消息元信息(代码生成专用)
func RegisterMessageMeta(name string, msgType reflect.Type, id uint32) {

	meta := &MessageMeta{
		Type: msgType,
		Name: name,
		ID:   id,
	}

	metaByName[name] = meta
	metaByID[meta.ID] = meta
}

// 根据名字查找消息元信息
func MessageMetaByName(name string) *MessageMeta {
	if v, ok := metaByName[name]; ok {
		return v
	}

	return nil
}

// 消息全名
func MessageFullName(rtype reflect.Type) string {

	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	return path.Base(rtype.PkgPath()) + "." + rtype.Name()

}

// 根据id查找消息元信息
func MessageMetaByID(id uint32) *MessageMeta {
	if v, ok := metaByID[id]; ok {
		return v
	}

	return nil
}

// 根据id查找消息名, 没找到返回空
func MessageNameByID(id uint32) string {

	if meta := MessageMetaByID(id); meta != nil {
		return meta.Name
	}

	return ""
}

// 遍历消息元信息
func VisitMessageMeta(callback func(*MessageMeta)) {

	for _, meta := range metaByName {
		callback(meta)
	}

}
